package namecheap

import (
	v1 "github.com/SabariVig/DNS-Operator/api/v1"
	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func convertToResponse(fetchedRecords *nc.DomainsDNSGetHostsCommandResponse) []nc.DomainsDNSHostRecord {
	var records []nc.DomainsDNSHostRecord
	for _, record := range *fetchedRecords.DomainDNSGetHostsResult.Hosts {
		records = append(records, nc.DomainsDNSHostRecord{
			HostName:   record.Name,
			RecordType: record.Type,
			Address:    record.Address,
		})
	}
	return records
}

func appendRecord(fetchedRecords *[]nc.DomainsDNSHostRecord, record *v1.Record) []nc.DomainsDNSHostRecord {
	return append(*fetchedRecords, nc.DomainsDNSHostRecord{
		Address:    &record.Spec.Address,
		HostName:   &record.Spec.Hostname,
		RecordType: (*string)(&record.Spec.Type),
	})

}
