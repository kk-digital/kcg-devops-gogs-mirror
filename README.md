# kcg-devops-gogs-mirror

## Document steps

## 1. Run gogs.

- https://github.com/gogs/gogs
- https://gogs.io/docs/installation/install_from_source.html

### Go for Gogs

#### Install from source

https://gogs.io/docs/installation/install_from_source.html

##### Installing Go

Gogs requires Go 1.18 to compile, please refer to the [official documentation](https://go.dev/doc/install) for how to install Go in your system.

```shell
$ go version
go version go1.19.8 linux/amd64
```

##### Set Up the Environment

We are going to create a new uand set up everything under that user:

```shell
sudo adduser --disabled-login --gecos 'Gogs' git

# CentOS:
sudo useradd -r -m -d /home/git -s /sbin/nologin -c "Gogs" git
```

##### Compile Gogs

###### Build with Tags

A couple of things do not come with Gogs automatically, you need to compile Gogs with corresponding build tags.

Available build tags are:

- pam: PAM authentication support
- cert: Generate self-signed certificates support
- minwinsvc: Builtin windows service support (or you can use NSSM to create a service)

```shell
# Clone the repository to the "gogs" subdirectory
git clone --depth 1 https://github.com/gogs/gogs.git gogs
# Change working directory
cd gogs
# Compile the main program, dependencies will be downloaded at this step
go build -tags "pam cert" -o gogs
```

If you get error: `fatal error: security/pam_appl.h: No such file or directory`, then install missing package via:

```shell
sudo apt-get install libpam0g-dev

# CentOS
sudo yum install -y pam-devel
```

##### Configure PAM policies

Edit or create a new PAM configuration file, such as `/etc/pam.d/gogs`, where gogs is the PAM service name you used in the Gogs configuration file. Add the following to the PAM configuration file:

```text
auth        required      pam_unix.so nullok
account     required      pam_unix.so
```

To configure this you just need to set the PAM Service Name to a filename in `/etc/pam.d/`. If you want it to work with normal Linux passwords, the user running Gogs must have read access to `/etc/shadow`.

##### Custom configuration

```shell
$ mkdir -p custom/conf/auth.d

$ vi custom/conf/auth.d/gogs_pam.conf
# Copy these contents(conf/gogs_pam.conf) into gogs_pam.conf
# This is an example of PAM authentication
#
id           = 104
type         = pam
name         = System Auth
is_activated = true

[config]
service_name = gogs

$ vi custom/conf/app.ini
# Copy these contents(conf/app.ini) into app.ini, replacing root with your current system user
BRAND_NAME = Gogs
RUN_USER   = root
RUN_MODE   = prod

[database]
TYPE     = sqlite3
HOST     = 127.0.0.1:5432
NAME     = gogs
SCHEMA   = public
USER     = gogs
PASSWORD =
SSL_MODE = disable
PATH     = data/gogs/data/gogs.db

[repository]
ROOT           = data/git/gogs-repositories
DEFAULT_BRANCH = master

[server]
DOMAIN           = localhost
HTTP_PORT        = 10880
EXTERNAL_URL     = http://localhost:10880/
DISABLE_SSH      = false
SSH_PORT         = 10022
START_SSH_SERVER = false
OFFLINE_MODE     = false

[mailer]
ENABLED = false

[auth]
REQUIRE_EMAIL_CONFIRMATION  = false
DISABLE_REGISTRATION        = false
ENABLE_REGISTRATION_CAPTCHA = true
REQUIRE_SIGNIN_VIEW         = false

[user]
ENABLE_EMAIL_NOTIFICATION = false

[picture]
DISABLE_GRAVATAR        = false
ENABLE_FEDERATED_AVATAR = false

[session]
PROVIDER = file

[log]
MODE      = file
LEVEL     = Info
ROOT_PATH = data/gogs/log

[security]
INSTALL_LOCK = true
SECRET_KEY   = d4wTSI1d3NRGzt0
```

##### Test Installation

To make sure Gogs is working:

```shell
./gogs web
```

If you do not see any error messages, hit Ctrl-C to stop Gogs.

## 2. Gogs may be run locally on any system with golang, or could be run in a docker container.

- like git pull, git clone of gogs

### Docker for Gogs

To keep your data out of Docker container, we do a volume ($HOME/gogs -> /data) here, and you can change it based on your situation.

```shell
# Pull image from Docker Hub.
$ docker pull gogs/gogs

# Create local directory for volume.
$ mkdir -p $HOME/gogs

# Use `docker run` for the first time.
$ docker run -d --name=gogs -p 10022:22 -p 10880:3000 -v $HOME/gogs:/data gogs/gogs

# Use `docker start` if you have stopped it.
$ docker start gogs
```

Note: It is important to map the SSH service from the container to the host and set the appropriate SSH Port and URI settings when setting up Gogs for the first time. To access and clone Git repositories with the above configuration you would use: `git clone ssh://git@hostname:10022/username/myrepo.git` for example.

In our example we have to type the following address into the browser:
`http://localhost:10880/install`

As a result you should see the setup page of Gogs:

<!-- ![Install Steps For First-time Run](docs/step1.png) -->

In this setup page we need to adapt the default settings to the settings we defined in the docker run command we executed previously. Thus, please change the input fields according to this table:

| Input field     | Description                                                      |
| --------------- | ---------------------------------------------------------------- |
| Database Type   | Replace PostgresSQL with SQLite3                                 |
| SSH Port        | Replace 22 with 10022                                            |
| HTTP Port       | Replace 3000 with 10880                                          |
| Application URL | Replace http://localhost:3000/ with http://localhost:10880/      |
| Log Path        | Replace /app/gogs/log with /data/gogs/log                        |
| Username        | Set an your username e.g. my-name                                |
| Password        | Set an your password with at least 8 characters e.g. my-Passw0rd |
| E-mail          | Set your email address e.g. demo@upwork.com                      |

After you filled all required fields, it should look like this:

<!-- ![Installed](docs/step2.png) -->

To finish the setup, click on **Install Gogs** at the bottom of the page.

After installing in you should see the start page of the Gogs service.

<!-- ![Start Page](docs/step3.png) -->

Nice, you successfully installed Gogs! Now let’s start using it! Please continue with the next section to learn how.

### How to create a Git repository in Gogs

Before we will be able to git push, we need to

- create a Git repository in Gogs
- configure your public SSH key in Gogs
- clone the Git repository in Gogs to your local workstation

Let’s start by creating a new repository. Click on the **blue button** in the top right corner of the page and choose **New Repository**.

<!-- ![New Repository](docs/step4.png) -->

You should see a setup page for your new repository. Please fill out the two text input fields.

At the bottom of the page you will find a checkbox. We recommend to enable this checkbox since it automatically will initialize this repository with a README.md. The following screenshot shows an example of how the setup page might look like after you provided all the necessary information.

<!-- ![Setup Page For New Repository](docs/step5.png) -->

Finish the setup by clicking on **Create Repository**.

Before you can actually use the new repository, you need to add your SSH key. Click on The settings button in the top right corner of the page. Then, go to **SSH Keys -> Add Key**.

<!-- ![Your Settings](docs/step6.png)
![Add SSH Key](docs/step7.png) -->

Paste your public SSH key into this `Content` field and set an arbitrary `Key Name`, e.g. ssh-rsa.

If you don’t know where to find your SSH key, execute the following command in your terminal.

```shell
cat ~/.ssh/id_rsa.pub
```

As you can see in the following screenshot, you should see your SSH key as result.

<!-- ![SSH Key](docs/step8.png) -->

Finally, back in the Gogs website, click the green button **Add Key** to add your public SSH key. Afterwards you should see a message confirming that the key has been added successfully.

<!-- ![Add Key](docs/step9.png)
![New SSH key 'ssh-rsa' has been added successfully!](docs/step10.png) -->

At this point, Gogs should be set up properly to receive your first commit! As a last step, we will `git clone` the new repo to your workstation.

#### How to generate gogs access token on command line

### Start using Git with the new repository

Gogs helps us in constructing the proper `git clone` command. Navigate inside your new repository in the Gogs Web-GUI and click on **SSH**. Make sure that the SSH button in front of the command is activated and copy the string in the text field via the **Copy** button on the right edge.

<!-- ![SSH Copy](docs/step11.png) -->

Now open a terminal at your workstation and navigate to the location where you want to create the folder for the repository.

Paste the command you just copied from Gogs Web-GUI. Before executing add git clone ssh:// in front of the command and put port 10022 in front of your user name. In our case 10022 is the port Gogs listens on for SSH.

In the end the command should look similar to this one:

```shell
git clone ssh://git@localhost:10022/my-name/demo-upwork.git
```

After this command has been executed confirm the fingerprint prompt and navigate into the new directory that has been created by this command.

Within the repository folder execute a `git pull`. If the command returns Already up-to-date the repository is properly set up.

<!-- ![Congratulations!! You now have your own Git service running on your workstation!](docs/step12.png) -->

## - how to create token

### Script - Create github token

```shell
$ go run main.go create-access-token -u root -p root -n script-token1
2023/04/23 17:04:56 Getting access token with gogs...
2023/04/23 17:04:56 user root created access token, Name: script-token1, Sha1: 0cbdfdd16430c7284aa39eb23ce3843b5dc8ef53
2023/04/23 17:04:56 Successfully getted acccess token, total cost: 33.973863ms
```

### GUI - Create github token

https://github.blog/2013-05-16-personal-api-tokens/

### Create gogs token

To obtain an API token for Gogs, you need to follow these steps:

1. Log in to your Gogs account.

2. Click on your avatar in the upper right corner of the webpage and then click **Your Settings**.

3. In the settings page, select **Applications** from the sidebar.
   <!-- ![Applications](docs/step13.png) -->

4. In the "Manage Personal Access Tokens" section, and click on **Generate New Token**.

5. Provide a descriptive "Token Name" for the token (for example, "My API token") and click **Generate Token**.
   <!-- ![Generate Token](docs/step14.png) -->

6. Gogs will generate a new API token for you. Be sure to copy the token to a safe place, as you won't be able to see it again. If you lose the token, you'll need to generate a new one.

You can now use this API token to authenticate your requests to the Gogs API. Remember to include the token in the "Authorization" header of your HTTP requests, like this:

## - how to run script

### build

```go
go build -o bin/gogs-helper

// or go run directly
go run main.go
```

### help

```go
./bin/gogs-helper --help
```

A helper tool to clone and update repositories between GitHub and Gogs

Usage:
gogs-helper [command]

Available Commands:
clone Clone all repos from GitHub organization to Gogs
completion Generate the autocompletion script for the specified shell
help Help about any command
update Update all existing repos in Gogs

Flags:
-t, --github-token string GitHub access token (default "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")
-s, --gogs-ssh-url string Gogs ssh URL (default "localhost:10022")
-g, --gogs-token string Gogs base URL (default "77cae12a2134d6e6ad8da5262a90502a412d7c03")
-u, --gogs-url string Gogs base URL (default "localhost:10880")
-n, --gogs-user-name string your Gogs user name (default "my-name")
-h, --help help for gogs-helper
-w, --workers int Speed up the command (default 6)

Use "gogs-helper [command] --help" for more information about a command.

### clone

Clone all repos from GitHub organization to Gogs

```go
go run main.go clone -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03

// or
./bin/gogs-helper clone -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03
```

### update

Update all existing repos in Gogs

```go
go run main.go update -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03

// or
./bin/gogs-helper update -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03
```

### list-org

Get a list of github organizations

```go
go run main.go list-org -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03

// or
./bin/gogs-helper list-org -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03
```

### list-org-repo

Get a list of github repositories in an organization

```go
go run main.go list-org-repo -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03

// or
./bin/gogs-helper list-org-repo -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03
```

### clone-local

Clone all repos from GitHub organization into a local directory

```go
go run main.go clone-local -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03

// or
./bin/gogs-helper clone-local -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03
```

### add

Add all repositories from a local directory to Gogs

```go
go run main.go add -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03

// or
./bin/gogs-helper add -t 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 -g 77cae12a2134d6e6ad8da5262a90502a412d7c03
```
