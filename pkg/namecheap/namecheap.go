package namecheap

import (
	"context"
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/SabariVig/DNS-Operator/api/v1"
	"github.com/SabariVig/DNS-Operator/pkg/providers"
	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	corev1 "k8s.io/api/core/v1"
)

type NCProvider struct {
	*nc.Client
}

func NewNamecheapProvider(secret corev1.Secret) providers.Providers {
	useSandbox, err := strconv.ParseBool(string(secret.Data["useSandbox"]))
	if err != nil {
		return nil
	}
	client := nc.NewClient(&nc.ClientOptions{
		UserName:   string(secret.Data["username"]),
		ApiUser:    string(secret.Data["apiUser"]),
		ApiKey:     string(secret.Data["apiKey"]),
		ClientIp:   string(secret.Data["clientIP"]),
		UseSandbox: useSandbox,
	})
	return &NCProvider{
		client,
	}
}

func (client *NCProvider) AddRecord(ctx context.Context, record *v1.Record) error {
	log := log.FromContext(ctx)
	fetchedRecords, err := client.DomainsDNS.GetHosts(record.Spec.Domain)
	if err != nil {
		log.Error(err, "Unable to get Hosts from Namecheap")
		return err
	}

	records := convertToResponse(fetchedRecords)
	log.Info("Add Record", "records", records)

	appendRecord(&records, record)

	log.Info("AddRecord: Setting Domain records", records)
	_, err = client.DomainsDNS.SetHosts(
		&nc.DomainsDNSSetHostsArgs{
			Domain:  &record.Spec.Domain,
			Records: &records,
		},
	)
	if err != nil {
		log.Error(err, "Error: Unable to Set Hosts from Namecheap")
		return err
	}
	log.Info("Add Record: Successfully set Domain for", "host_name", record.Spec.Hostname, "domain_name", record.Spec.Domain)
	return nil
}

func (client *NCProvider) DeleteRecord(ctx context.Context, record *v1.Record) error {
	log := log.FromContext(ctx)
	log.Info("DeleteRecord: started process")

	fetchedRecords, err := client.DomainsDNS.GetHosts(record.Spec.Domain)
	if err != nil {
		log.Error(err, "Unable to get Hosts from Namecheap")
		return err
	}

	records := convertToResponse(fetchedRecords)
	log.Info("DeleteRecord", "Exsisting Record", records)
	filteredRecord := removeRecord(records, record)

	if len(records) == len(filteredRecord) {
		log.Info("No record to delete")
		return nil
	}

	_, err = client.DomainsDNS.SetHosts(
		&nc.DomainsDNSSetHostsArgs{
			Domain:  &record.Spec.Domain,
			Records: &filteredRecord,
		},
	)

	if err != nil {
		log.Error(err, "DeleteRecord: Unable to delete from Namecheap")
		return err
	}

	log.Info("Successfully Deleted Record from Namecheap for", record.Spec.Domain, record.Spec.Domain)
	return nil
}
