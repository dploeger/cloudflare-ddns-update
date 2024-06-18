// Package pkg includes the publicly available API for cloudflare-ddns-update
package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go/v2"
	"github.com/cloudflare/cloudflare-go/v2/dns"
	"github.com/cloudflare/cloudflare-go/v2/option"
	"github.com/cloudflare/cloudflare-go/v2/zones"
	"github.com/dploeger/cloudflare-ddns-update/internal"
	"github.com/gin-gonic/gin"
	"net/netip"
)

// The API includes all available endpoints
type API struct {
	// cfClient is the currently used cloudflare client
	cfClient *cloudflare.Client
	// zoneID is the id where the records are managed in
	zoneID string
}

// NewAPI creates a new API variable. apiToken is used for authentication to cloudflare and debug enables the
// request/response logger
func NewAPI(apiToken string, zoneName string, debug bool) (*API, error) {
	options := []option.RequestOption{option.WithAPIToken(apiToken)}
	if debug {
		options = append(options, option.WithMiddleware(internal.RequestResponseLogger))
	}
	a := API{
		cfClient: cloudflare.NewClient(options...),
	}

	z := a.cfClient.Zones.ListAutoPaging(context.Background(), zones.ZoneListParams{Name: cloudflare.F(zoneName)})
	if z.Err() != nil {
		return &a, fmt.Errorf("error getting zone ID for zone %s: %s", zoneName, z.Err())
	}
	if !z.Next() {
		return &a, fmt.Errorf("can't get zone ID for zone %s", zoneName)
	}
	a.zoneID = z.Current().ID

	return &a, nil
}

// DDNSUpdate is based on the Dyn v3 API and should be bound to GET /v3/update
func (a API) DDNSUpdate(c *gin.Context) {
	var params struct {
		Hostname string `form:"hostname"`
		MyIP     string `form:"myip"`
	}
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Missing query parameters: %s", err),
		})
		return
	}
	var myIP netip.Addr
	if a, err := netip.ParseAddr(params.MyIP); err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Error converting IP address: %s", err),
		})
		return
	} else {
		myIP = a
	}

	var recordType string
	if r, ok := c.Params.Get("type"); !ok {
		if myIP.Is6() {
			recordType = "AAAA"
		} else {
			recordType = "A"
		}
	} else {
		recordType = r
	}

	records := a.cfClient.DNS.Records.ListAutoPaging(context.Background(), dns.RecordListParams{Name: cloudflare.F(params.Hostname), ZoneID: cloudflare.F(a.zoneID)})

	if records.Next() {
		updateParams := dns.RecordUpdateParams{
			ZoneID: cloudflare.F(a.zoneID),
			Record: dns.ARecordParam{
				Name:    cloudflare.F(records.Current().Name),
				Type:    cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F(myIP.String()),
			},
		}
		if _, err := a.cfClient.DNS.Records.Update(context.Background(), records.Current().ID, updateParams); err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("Can't update DNS record for host %s from zone %s: %s", params.Hostname, a.zoneID, err),
			})
			return
		}
	} else {
		if recordType == "A" {
			if _, err := a.cfClient.DNS.Records.New(context.Background(), dns.RecordNewParams{
				ZoneID: cloudflare.F(a.zoneID),
				Record: dns.ARecordParam{
					Type:    cloudflare.F(dns.ARecordTypeA),
					Name:    cloudflare.F(params.Hostname),
					Content: cloudflare.F(myIP.String()),
				},
			}); err != nil {
				c.JSON(500, gin.H{
					"error": fmt.Sprintf("Can't create DNS record for host %s in zone %s: %s", params.Hostname, a.zoneID, err),
				})
				return
			}
		}
	}
	c.Status(204)
}
