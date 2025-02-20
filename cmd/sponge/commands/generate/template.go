package generate

const (
	mainFileHTTPCode = `
	// 创建http服务
	httpAddr := ":" + strconv.Itoa(config.Get().HTTP.Port)
	httpRegistry, httpInstance := registryService("http", config.Get().App.Host, config.Get().HTTP.Port)
	httpServer := server.NewHTTPServer(httpAddr,
		server.WithHTTPReadTimeout(time.Second*time.Duration(config.Get().HTTP.ReadTimeout)),
		server.WithHTTPWriteTimeout(time.Second*time.Duration(config.Get().HTTP.WriteTimeout)),
		server.WithHTTPRegistry(httpRegistry, httpInstance),
		server.WithHTTPIsProd(config.Get().App.Env == "prod"),
	)
	servers = append(servers, httpServer)
`

	mainFileGrpcCode = `
	// 创建grpc服务
	grpcAddr := ":" + strconv.Itoa(config.Get().Grpc.Port)
	grpcRegistry, grpcInstance := registryService("grpc", config.Get().App.Host, config.Get().Grpc.Port)
	grpcServer := server.NewGRPCServer(grpcAddr,
		server.WithGrpcReadTimeout(time.Duration(config.Get().Grpc.ReadTimeout)*time.Second),
		server.WithGrpcWriteTimeout(time.Duration(config.Get().Grpc.WriteTimeout)*time.Second),
		server.WithGrpcRegistry(grpcRegistry, grpcInstance),
	)
	servers = append(servers, grpcServer)
`

	dockerFileHTTPCode = `# 添加curl，用在http服务的检查，如果用部署在k8s，可以不用安装
RUN apk add curl

COPY configs/ /app/configs/
COPY serverNameExample /app/serverNameExample
RUN chmod +x /app/serverNameExample

# http端口
EXPOSE 8080`

	dockerFileGrpcCode = `# 添加grpc_health_probe，用在grpc服务的健康检查
COPY grpc_health_probe /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe

COPY configs/ /app/configs/
COPY serverNameExample /app/serverNameExample
RUN chmod +x /app/serverNameExample`

	dockerFileBuildHTTPCode = `# 添加curl，用在http服务的检查，如果用部署在k8s，可以不用安装
RUN apk add curl

COPY --from=build /serverNameExample /app/serverNameExample
COPY --from=build /go/src/serverNameExample/configs/serverNameExample.yml /app/configs/serverNameExample.yml

# http端口
EXPOSE 8080`

	dockerFileBuildGrpcCode = `# 添加grpc_health_probe，用在grpc服务的健康检查
COPY --from=build /grpc_health_probe /bin/grpc_health_probe
COPY --from=build /serverNameExample /app/serverNameExample
COPY --from=build /go/src/serverNameExample/configs/serverNameExample.yml /app/configs/serverNameExample.yml`

	dockerComposeFileHTTPCode = `    ports:
      - "8080:8080"   # http端口
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]   # http健康检查，注：镜像必须包含curl命令`

	dockerComposeFileGrpcCode = `
    ports:
      - "8282:8282"   # grpc服务端口
      - "8283:8283"   # grpc metrics端口
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:8282"]    # grpc健康检查，注：镜像必须包含grpc_health_probe命令`

	k8sDeploymentFileHTTPCode = `
          ports:
            - name: http-port
              containerPort: 8080
          readinessProbe:
            httpGet:
              port: http-port
              path: /health
            initialDelaySeconds: 10
            timeoutSeconds: 2
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          livenessProbe:
            httpGet:
              port: http-port
              path: /health`

	k8sDeploymentFileGrpcCode = `
          ports:
            - name: grpc-port
              containerPort: 8282
            - name: metrics-port
              containerPort: 8283
          readinessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:8282"]
            initialDelaySeconds: 10
            timeoutSeconds: 2
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          livenessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:8282"]`

	k8sServiceFileHTTPCode = `  ports:
    - name: server-name-example-svc-http-port
      port: 8080
      targetPort: 8080`

	k8sServiceFileGrpcCode = `  ports:
    - name: server-name-example-svc-grpc-port
      port: 8282
      targetPort: 8282
    - name: server-name-example-svc-grpc-metrics-port
      port: 8283
      targetPort: 8283`

	configFileCode = `// nolint
// code generated by sponge.

package config

import (
	"github.com/zhufuyi/sponge/pkg/conf"
)

var config *Config

func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Show() string {
	return conf.Show(config)
}

func Get() *Config {
	if config == nil {
		panic("config is nil")
	}
	return config
}

func Set(conf *Config) {
	config = conf
}
`

	configFileCcCode = `// nolint
// code generated by sponge.

package config

import (
	"github.com/zhufuyi/sponge/pkg/conf"
)

func NewCenter(configFile string) (*Center, error) {
	nacosConf := &Center{}
	err := conf.Parse(configFile, nacosConf)
	return nacosConf, err
}
`
)
