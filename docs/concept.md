# Konzept

Borg ist ein Server der eine minimale REST-API bereitstellt um beliebige Dateien zu analysieren. Die Analyse besteht aus Formaterkennung, Metadatenextraktion und Formatvalidierung.

## Motivation

Die Formaterkennung und -validierung von unbekannten Dateien ist ein komplexes Problem. Aufgrund der Komplexität kann kein einzelnes Programm das Problemfeld vollständig lösen. In der Regel sind Anwendungen darauf spezialisiert, entweder Dateien mit unbekannten Formaten zu identifizieren oder eine Formatfamilie bzw. eine kleine Gruppe von Dateiformaten zu validieren.
Um eine möglichst umfassende Abdeckung bei der Identifizierung und Validierung von Dateiformaten zu erreichen, ist es daher notwendig, mehrere Programme miteinander zu kombinieren. Es gibt bereits vereinzelt Programme, die mehrere Werkzeuge für die Formaterkennung und -validierung einbinden. Diese eingebundenen Werkzeuge werden in der Regel direkt integriert oder lokal ausgeführt.
Für Borg wurde ein anderer Ansatz gewählt. Die Werkzeuge werden nicht direkt integriert, sondern werden in eigenen Containern ausgeführt und über eine Web-API angesprochen. Das verringert die Abhängigkeit von Systemvorraussetzungen der verwendeten Werkzeuge.

## Werkzeugserver

Die integrierten Werkzeuge werden containerisiert betrieben.

## Eigenschaftsmengen

_PUID_, _MIME-Type_, _Formatversion_ und _Validität_ sind die zentralen Eigenschaften zur Bestimmung des Dateityps.
