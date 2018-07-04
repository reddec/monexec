package plugins

import (
	"context"
	"github.com/reddec/monexec/pool"
	"github.com/pkg/errors"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"os"
	"io"
	"github.com/elazarl/go-bindata-assetfs"
	"fmt"
)

const restApiStartupCheck = 1 * time.Second

//go:generate go-bindata -pkg plugins -prefix ../ui/dist/ ../ui/dist/
type RestPlugin struct {
	Listen string `yaml:"listen"`
	server *http.Server
}

func (p *RestPlugin) Prepare(ctx context.Context, pl *pool.Pool) error {

	router := gin.Default()
	router.StaticFS("/ui/", &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: ""})
	router.GET("/supervisors", func(gctx *gin.Context) {
		var names = make([]string, 0)
		for _, sv := range pl.Supervisors() {
			names = append(names, sv.Config().Name)
		}
		gctx.JSON(http.StatusOK, names)
	})
	router.GET("/supervisor/:name", func(gctx *gin.Context) {
		name := gctx.Param("name")
		for _, sv := range pl.Supervisors() {
			if sv.Config().Name == name {
				gctx.JSON(http.StatusOK, sv.Config())
				return
			}
		}
		gctx.AbortWithStatus(http.StatusNotFound)
	})
	router.GET("/supervisor/:name/log", func(gctx *gin.Context) {
		name := gctx.Param("name")
		for _, sv := range pl.Supervisors() {
			if sv.Config().Name == name {
				if sv.Config().LogFile == "" {
					break
				}
				f, err := os.Open(sv.Config().LogFile)
				if err != nil {
					gctx.AbortWithError(http.StatusBadGateway, err)
					return
				}
				defer f.Close()
				gctx.Header("Content-Type", "text/plain")
				gctx.Header("Content-Disposition", "attachment; filename=\""+sv.Config().Name+".log\"")
				gctx.AbortWithStatus(http.StatusOK)
				io.Copy(gctx.Writer, f)
				return
			}
		}
		gctx.AbortWithStatus(http.StatusNotFound)
	})
	router.POST("/supervisor/:name", func(gctx *gin.Context) {
		name := gctx.Param("name")
		for _, sv := range pl.Supervisors() {
			if sv.Config().Name == name {
				in := pl.Start(ctx, sv)
				gctx.JSON(http.StatusOK, in)
				return
			}
		}
		gctx.AbortWithStatus(http.StatusNotFound)
	})
	router.GET("/instances", func(gctx *gin.Context) {
		var names = make([]string, 0)
		for _, sv := range pl.Instances() {
			names = append(names, sv.Config().Name)
		}
		gctx.JSON(http.StatusOK, names)
	})

	router.GET("/instance/:name", func(gctx *gin.Context) {
		name := gctx.Param("name")
		for _, sv := range pl.Instances() {
			if sv.Config().Name == name {
				gctx.JSON(http.StatusOK, sv)
				return
			}
		}
		gctx.AbortWithStatus(http.StatusNotFound)
	})

	router.POST("/instance/:name", func(gctx *gin.Context) {
		name := gctx.Param("name")
		for _, sv := range pl.Instances() {
			if sv.Config().Name == name {
				pl.Stop(sv)
				gctx.AbortWithStatus(http.StatusCreated)
				return
			}
		}
		gctx.AbortWithStatus(http.StatusNotFound)
	})

	p.server = &http.Server{Addr: p.Listen, Handler: router}
	fmt.Println("rest interface will be available on", p.Listen)
	start := make(chan error, 1)
	go func() {
		start <- p.server.ListenAndServe()
	}()
	select {
	case err := <-start:
		return err
	case <-time.After(restApiStartupCheck):
		return nil
	}
}

func (p *RestPlugin) OnSpawned(ctx context.Context, sv pool.Instance) {}

func (p *RestPlugin) OnStarted(ctx context.Context, sv pool.Instance) {}

func (p *RestPlugin) OnStopped(ctx context.Context, sv pool.Instance, err error) {}

func (p *RestPlugin) OnFinished(ctx context.Context, sv pool.Instance) {}

func (p *RestPlugin) MergeFrom(o interface{}) error {
	def := defaultRestPlugin()
	other := o.(*RestPlugin)
	if p.Listen == def.Listen {
		p.Listen = other.Listen
	} else if other.Listen != def.Listen && other.Listen != p.Listen {
		return errors.Errorf("unmatched Rest listen address %v != %v", p.Listen, other.Listen)
	}
	return nil
}

func (a *RestPlugin) Close() error {
	ctx, closer := context.WithTimeout(context.Background(), 1*time.Second)
	defer closer()
	return a.server.Shutdown(ctx)
}

func defaultRestPlugin() *RestPlugin {
	return &RestPlugin{
		Listen: "localhost:9900",
	}
}

func init() {
	registerPlugin("rest", func(file string) PluginConfigNG {
		return defaultRestPlugin()
	})
}
