package main

import (
	"goapp/cm"
	"goapp/pb"
	"goapp/sv"
	"goapp/sy"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	userServiceClient := pb.NewUserServiceClient(conn)

	appService := sv.AppService{UserClient: userServiceClient}

	app.Route("/", func() app.Composer { return &cm.Dashboard{} })
	app.Route("/login", func() app.Composer { return &sy.Login{AppService: appService} })
	app.Route("/sy", func() app.Composer { return &sy.Menu{} })

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Rbooks",
		Description: "Yet another quickbooks like app",
		//Resources:   app.LocalDir("/app/web"),
		//Resources:   app.RemoteBucket("https://storage.googleapis.com/myapp.appspot.com"),
		Icon: app.Icon{
			Default: "/web/images/favicon.ico",
			Large:   "/web/images/logo.png",
		},
		Styles: []string{
			"/web/css/app.css", // Loads hello.css file.
			"/web/css/bulma.css",
		},
		Scripts: []string{
			"/web/js/app.js", // Loads hello.js file.
			// "https://foo.com/remoteScript.js", // loads a remote script.
		},
		RawHeaders: []string{
			`<link rel="icon" type="image/x-icon" href="/web/images/favicon.ico">`,
			`<link rel="manifest" href="/web/manifest.webmanifest">`,
		},
	})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
	log.Fatal(http.ListenAndServe(":8000", nil))

}
