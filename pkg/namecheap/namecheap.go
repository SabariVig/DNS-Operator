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
	fechRecords := convertToResponse(fetchedRecords)
	log.Info("Exsisting Record", fechRecords)

	records := appendRecord(&fechRecords, record)

	log.Info("Setting Domain records", records)
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
	return nil
}

func (client *NCProvider) DeleteRecord(ctx context.Context, record *v1.Record) error {
	return nil
}
