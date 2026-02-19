module github.com/antinvestor/service-ocr

go 1.25.2

require (
	buf.build/gen/go/antinvestor/common/protocolbuffers/go v1.36.11-20260219054105-fe125014d75c.1
	buf.build/gen/go/antinvestor/files/connectrpc/go v1.19.1-20251223175512-7a3fcba332ca.2
	buf.build/gen/go/antinvestor/files/protocolbuffers/go v1.36.11-20251223175512-7a3fcba332ca.1
	buf.build/gen/go/antinvestor/ocr/connectrpc/go v1.19.1-20260102123348-6bcf22d3ced1.2
	buf.build/gen/go/antinvestor/ocr/protocolbuffers/go v1.36.11-20260102123348-6bcf22d3ced1.1
	connectrpc.com/connect v1.19.1
	github.com/antinvestor/apis/go/common v1.53.5
	github.com/antinvestor/apis/go/files v1.53.6
	github.com/otiai10/gosseract/v2 v2.4.1
	github.com/pitabwire/frame v1.72.1
	github.com/pitabwire/util v0.4.0
	golang.org/x/oauth2 v0.35.0
	google.golang.org/api v0.266.0
	google.golang.org/protobuf v1.36.11
	gorm.io/datatypes v1.2.7
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1 // indirect
	buf.build/gen/go/gnostic/gnostic/protocolbuffers/go v1.36.11-20230414000709-087bc8072ce4.1 // indirect
	buf.build/go/protovalidate v1.1.0 // indirect
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go/auth v0.18.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	connectrpc.com/otelconnect v0.9.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/caarlos0/env/v11 v11.3.1 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/exaring/otelpgx v0.10.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/cel-go v0.27.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.11 // indirect
	github.com/googleapis/gax-go/v2 v2.17.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/lmittmann/tint v1.1.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nats.go v1.48.0 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.6.1 // indirect
	github.com/ory/keto/proto v0.13.0-alpha.0 // indirect
	github.com/panjf2000/ants/v2 v2.11.5 // indirect
	github.com/pitabwire/natspubsub v0.7.12 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/sony/gobreaker/v2 v2.4.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/bridges/otelslog v0.15.0 // indirect
	go.opentelemetry.io/contrib/bridges/prometheus v0.65.0 // indirect
	go.opentelemetry.io/contrib/exporters/autoexport v0.65.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.65.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.65.0 // indirect
	go.opentelemetry.io/contrib/propagators/autoprop v0.65.0 // indirect
	go.opentelemetry.io/contrib/propagators/aws v1.40.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.40.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.40.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.40.0 // indirect
	go.opentelemetry.io/otel v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.62.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.16.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.40.0 // indirect
	go.opentelemetry.io/otel/log v0.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/sdk v1.40.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.16.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/trace v1.40.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	gocloud.dev v0.44.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/exp v0.0.0-20260112195511-716be5621a96 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260203192932-546029d2fa20 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260203192932-546029d2fa20 // indirect
	google.golang.org/grpc v1.78.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
	gorm.io/gorm v1.31.1 // indirect
)
