// Copyright 2019, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lightstepexporter

import (
	"context"

	lightstepoc "github.com/lightstep/lightstep-tracer-go/lightstepoc"
	"github.com/spf13/viper"

	"github.com/census-instrumentation/opencensus-service/consumer"
	"github.com/census-instrumentation/opencensus-service/exporter/exporterwrapper"
)

type lightStepConfig struct {
	SatelliteHost         string `mapstructure:satellite_host,omitempty"`
	SatellitePort         int    `mapstructure:satellite_port,omitempty"`
	SatelliteUsePlaintext bool   `mapstructure:satellite_use_plaintext,omitempty"`
	AccessToken           string `mapstructure:access_token,omitempty"`
}

// LightStepExportersFromViper unmarshals the viper and returns exporter.TraceExporters targeting LightStep per config settings
func LightStepExportersFromViper(v *viper.Viper) (tps []consumer.TraceConsumer, mps []consumer.MetricsConsumer, doneFns []func() error, err error) {
	var config struct {
		LightStep *lightStepConfig `mapstructure:"lightstep"`
	}
	if err := v.Unmarshal(&config); err != nil {
		return nil, nil, nil, err
	}
	lsConfig := config.LightStep
	if lsConfig == nil {
		return nil, nil, nil, nil
	}

	lsExporter, err := lightstepoc.NewExporter(lightstepoc.WithInsecure(lsConfig.SatelliteUsePlaintext))
	if err != nil {
		return nil, nil, nil, err
	}

	doneFns = append(doneFns, func() error {
		lsExporter.Flush(context.Background())
		return nil
	})

	lsTraceExporter, err := exporterwrapper.NewExporterWrapper("lightstep", "ocservice.exporter.LightStep.ConsumeTraceData", lsExporter)
	if err != nil {
		return nil, nil, nil, err
	}

	tps = append(tps, lsTraceExporter)
	return
}
