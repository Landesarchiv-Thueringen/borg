# Installation

Für den Betrieb von Borg wird [Docker](https://docs.docker.com/) inklusive [Docker Compose](https://docs.docker.com/compose/) benötigt. Für den regulären Betrieb empfehlen wir die Installation auf einem Linux-Server. Für einen lokalen Test der Standalone-Version von Borg ist auch der Einsatz von [Docker Desktop](https://docs.docker.com/desktop/) möglich. In der Datei [.env](.env) kann ein Proxy für den Zugriff auf das Internet und der Port unter dem die Anwendung angesprochen werden soll, festgelegt werden. Borg benötigt nur bei der Erstellung der Container eine Internetverbindung. Im Betrieb wird keine Internetverbindung benötigt.

Um die Anwendung in einem Netzwerk verfügbar zu machen, eignet sich ein Webserver als Reverse-Proxy wie bspw. [nginx](https://www.nginx.com/), der die Anfragen auf den konfigurierten Port der Anwendung weiterleitet. Bei der Konfiguration des Webservers ist zu beachten, dass die Grenzen für die Dateigröße beim Upload und Timeouts von Verbindungen ausreichend bemessen werden, so dass auch größere Dateien analysiert werden können. Die notwendigen Einstellungen sind für alle Webserver unterschiedlich.

Folgend ist eine mögliche Teilkonfiguration für nginx abgebildet, die Dateigröße und Timeouts für den Upload erhöht:

**Beispiel nginx**

```nginx
location / {
    proxy_pass http://url-to-service/;

    client_max_body_size 5000m;
    proxy_connect_timeout 600;
    proxy_send_timeout 600;
    proxy_read_timeout 600;
    send_timeout 600;
}
```

Um alle für den Betrieb von Borg benötigten Container zu starten, genügen die Befehle:

```sh
cp .env.example .env
docker compose up --build -d
```
