module watcher

replace gitlab.sas.com/async-event-infrastructure/server => ../../../

go 1.21.0

require gitlab.sas.com/async-event-infrastructure/server v0.0.0-20230901192240-3d330b648418

require (
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/twmb/franz-go v1.14.4 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.6.1 // indirect
)
