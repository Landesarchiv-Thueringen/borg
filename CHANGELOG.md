# Changelog

## Next

- Dokumentation: SBOM für Container-Images bereitgestellt

## 2.1.0

- Fix: CORS für Integration in andere Anwendungen wieder erlaubt
- Intern: Timeouts für alle Werkzeuge
- Intern/Integration: Umbenennung Hauptendpunkt

## 2.0.0

- Feature: neues Konzept zur Berechnung des Gesamtergebnisses
- Feature: neue Ansicht für Metadaten
- Feature: Update DROID 6.8.0 -> 6.8.1
- Feature: Integration Siegfried 1.11.2
- Feature: Integration MediaInfo 24.11
- Dokumentation: neue [Dokumentation](https://landesarchiv-thueringen.github.io/borg)

## 1.4.1

- Feature: Version in Titelleiste und "Über Borg"-Dialog
- Fix: Antwort auf `/api/version`
- Fix: Schreibfehler in angezeigten Nachrichten
- Intern: Skript für das Veröffentlichen von Docker-Images

## 1.4.0

- Feature: Integration Magika v0.6.0
- Fix: DROID-Fehler bei nicht erkannten Format
- Fix: Filter-Chips können ausgewählt werden
- Fix: HTML-Attribut `lang`
- Fix: Darstellung boolescher Werte in Dateiübersicht
- Fix: Ersetzung falscher MIME-Types für Tika und Magika
- Fix: Darstellungsfehler bei langen Dateinamen
- Intern: Konfliktauflösung Formatversion
- Intern: Angular-Update auf Version 19
- Intern: Unterstützung für verschiedene Typen extrahierter Werte

## 1.3.0

- Fix: Fehler bei erster Ausführung von DROID
- Intern: Erlaube Konfiguration zur Laufzeit
- Intern: Konfigurierbare Ports
- Intern: Konfigurierbare Docker-Image-Namen
- Intern: Docker-Images optimieren
- Intern: Migration auf PNPM

## 1.2.0

- Feature: Update DROID 6.7.0 -> 6.8.0
- Feature: Update Tika 2.9.0 -> 2.9.2
- Feature: Update JHove 1.28.0 -> 1.30.1
- Feature: Update veraPDF 1.24.1 -> 1.26.2
- Feature: Update OOXML Validator 2.1.1 -> 2.1.5
- Fix: Reverse-Proxy-Weiterleitung von API-Requests

## 1.1.0

- Feature: Drag-and-drop-Unterstützung
- Feature: Umstieg auf Material-Design 3
- Feature: Sortierbare Ergebnis-Tabelle
- Feature: Filterbare Ergebnis-Tabelle
- Feature: Verbesserte Fehlerbehandlung
- Feature: Verschiedene kleinere UI-Anpassungen
- Feature: Häufiger automatische Umleitung zu passender Seite
- Intern: Ergebnisdatenformat überarbeitet (API)
- Intern: Endpunkt zum Prüfen der Borg-Version (API)
- Intern: Angular-Update von Version 16 auf 18
- Intern: Modulloses Angular-Projekt-Layout
- Intern: Verbesserungen bei Konfiguration
