package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/paked/messenger"
)

//profileField is a slice of strings of the user profile field the developer wants access
var (
	profileField = []string{"name", "first_name", "last_name", "profile_pic"}
)

const (
	webhooksPath = "/webhooks"
	loginPath    = "/signin"
	logoutPath   = "/signout"

	validUsername = "john"
	validPassword = "secret"
)

var (
	verifyToken = flag.String("verify-token", "", "The token used to verify facebook (required)")
	pageToken   = flag.String("page-token", "", "The token that is used to verify the page on facebook.")
	appSecret   = flag.String("app-secret", "", "The app secret from the facebook developer portal (required)")
	host        = flag.String("host", "localhost", "The host used to serve the messenger bot")
	port        = flag.Int("port", 8080, "The port used to serve the messenger bot")
	publicHost  = flag.String("public-host", "example.org", "The public facing host used to access the messenger bot")
)

func main() {
	flag.Parse()

	if *verifyToken == "" || *appSecret == "" || *pageToken == "" {
		fmt.Println("missing arguments")
		fmt.Println()
		flag.Usage()

		os.Exit(-1)
	}

	// Instantiate messenger client
	client := messenger.New(messenger.Options{
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	// Handle incoming messages
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		log.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		p, err := client.ProfileByID(m.Sender.ID, profileField)
		if err != nil {
			log.Println("Failed to fetch user profile:", err)
		}

		switch strings.ToLower(m.Text) {
		case "login":
			err = loginButton(r)
		case "logout":
			err = logoutButton(r)
		case "help":
			err = help(p, r)
		default:
			err = greeting(p, r)
		}

		if err != nil {
			log.Println("Failed to respond:", err)
		}
	})

	// Send a feedback to the user after an update of account linking status
	client.HandleAccountLinking(func(m messenger.AccountLinking, r *messenger.Response) {
		var text string
		switch m.Status {
		case "linked":
			text = "Hey there! You're now logged in :)"
		case "unlinked":
			text = "You've been logged out of your account."
		}

		if err := r.Text(text, messenger.ResponseType); err != nil {
			log.Println("Failed to send account linking feedback")
		}
	})

	// Setup router
	mux := http.NewServeMux()
	mux.Handle(webhooksPath, client.Handler())
	mux.HandleFunc(loginPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			loginForm(w, r)
		case "POST":
			login(w, r)
		}
	})

	// Listen
	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Serving messenger bot on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// loginButton will present to the user a button that can be used to
// start the account linking process.
func loginButton(r *messenger.Response) error {
	buttons := &[]messenger.StructuredMessageButton{
		{
			Type: "account_link",
			URL:  "https://" + path.Join(*publicHost, loginPath),
		},
	}
	return r.ButtonTemplate("Link your account.", buttons, messenger.ResponseType)
}

// logoutButton show to the user a button that can be used to start
// the process of unlinking an account.
func logoutButton(r *messenger.Response) error {
	buttons := &[]messenger.StructuredMessageButton{
		{
			Type: "account_unlink",
		},
	}
	return r.ButtonTemplate("Unlink your account.", buttons, messenger.ResponseType)
}

// greeting salutes the user.
func greeting(p messenger.Profile, r *messenger.Response) error {
	return r.Text(fmt.Sprintf("Hello, %v!", p.FirstName), messenger.ResponseType)
}

// help displays possibles actions to the user.
func help(p messenger.Profile, r *messenger.Response) error {
	text := fmt.Sprintf(
		"%s, looking for actions to do? Here is what I understand.",
		p.FirstName,
	)

	replies := []messenger.QuickReply{
		{
			ContentType: "text",
			Title:       "Login",
		},
		{
			ContentType: "text",
			Title:       "Logout",
		},
	}

	return r.TextWithReplies(text, replies, messenger.ResponseType)
}

// loginForm is the endpoint responsible to displays a login
// form. During the account linking process, after clicking on the
// login button, users are directed to this form where they are
// supposed to sign into their account. When the form is submitted,
// credentials are sent to the login endpoint.
func loginForm(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	linkingToken := values.Get("account_linking_token")
	redirectURI := values.Get("redirect_uri")
	fmt.Fprint(w, templateLogin(loginPath, linkingToken, redirectURI, false))
}

// login is the endpoint that handles the actual signing in, by
// checking the credentials, then redirecting to Facebook Messenger if
// they are valid.
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.FormValue("username")
	password := r.FormValue("password")
	linkingToken := r.FormValue("account_linking_token")
	rawRedirect := r.FormValue("redirect_uri")

	if !checkCredentials(username, password) {
		fmt.Fprint(w, templateLogin(loginPath, linkingToken, rawRedirect, true))
		return
	}

	redirectURL, err := url.Parse(rawRedirect)
	if err != nil {
		log.Println("failed to parse url:", err)
		return
	}

	q := redirectURL.Query()
	q.Set("authorization_code", "something")
	redirectURL.RawQuery = q.Encode()

	w.Header().Set("Location", redirectURL.String())
	w.WriteHeader(http.StatusFound)
}

func checkCredentials(username, password string) bool {
	return username == validUsername && password == validPassword
}

// templateLogin constructs the signin form.
func templateLogin(loginPath, linkingToken, redirectURI string, failed bool) string {
	failedInfo := ""
	if failed {
		failedInfo = `<p class="alert alert-danger">Incorrect credentials</p>`
	}

	template := `
<html>
  <head>
    <link rel="stylesheet"
          href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css">
  </head>
  <body>
    <div class="container">
    <h1>Login to your account</h1>
    <p class="alert alert-primary">
      Valid credentials are "%s" as the username and "%s" as the password
    </p>
    %s
    <form action="%s" method="POST">
      <input type="text" name="username" placeholder="Username">
      <input type="password" name="password" placeholder="Password">
      <input type="hidden" name="account_linking_token" value="%s">
      <input type="hidden" name="redirect_uri" value="%s">
      <button type="submit">Submit</button>
    </form>
    </div>
  </body>
</html>
`
	return fmt.Sprintf(
		template,
		validUsername,
		validPassword,
		failedInfo,
		loginPath,
		linkingToken,
		redirectURI,
	)
}
