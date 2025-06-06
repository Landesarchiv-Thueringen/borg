# Benutzerhandbuch

Borg stellt eine Standalone-Webanwendung bereit, mit der beliebige Dateien analysiert werden können.

## Datei-Auswahl

Der Reiter Dateiauswahl ermöglicht die Auswahl von einzelnen Dateien und ganzen Ordnern für die Analyse. Die Auswahl kann per Dialog oder Drag and Drop erfolgen. Wenn ein Ordner ausgewählt wird, werden auch die Dateien aller enthaltenen Ordner hochgeladen.

<figure markdown="span">
  ![Dateiauswahl](img/upload_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 1: Datei-Auswahl</figcaption></center>
</figure>

## Auswertung

Im Reiter Auswertung wird das Gesamtergebnis für alle hochgeladenen Dateien dargestellt. Für jede Datei werden in einer eigenen Zeile die wichtigsten extrahierten Eigenschaften sowie der Status angezeigt. Der Status stellt die Qualität des Gesamtergebnisses symbolisch dar. Detaillierte Ergebnisse können durch einen Klick auf die Zeile aufgerufen werden.

<figure markdown="span">
  ![Auswertung](img/overview_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 2: Auswertung für alle Dateiergebnisse</figcaption></center>
</figure>

## Detailansicht von Dateiergebnissen

Die Detailansicht von Dateiergebnissen ist in zwei Reiter aufgeteilt. Im Reiter Überblick wird der Status erläutert und alle erkannten Dateiformate aufgelistet.

<figure markdown="span">
  ![Detailansicht Datei 1](img/file_details_1_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 3: Detailansicht einer validen JPEG-Datei</figcaption></center>
</figure>

Wie sich die Eigenschaften von den erkannten Dateiformaten zusammensetzen, wird in einer Detailansicht klar.

<figure markdown="span">
  ![Detailansicht Datei 1](img/format_result_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 4: Detailansicht für ein erkanntes Dateiformat</figcaption></center>
</figure>

Wenn mehrere mögliche Dateiformate erkannt werden, hilft die Bewertung bei der Einordnung der Ergebnisse.

<figure markdown="span">
  ![Detailansicht Datei 2](img/file_details_2_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 5: Detailansicht einer Datei mit mehreren erkannten Dateiformaten</figcaption></center>
</figure>

Werden Probleme mit einer Datei festgestellt wird das klar in der Oberfläche kommuniziert.

<figure markdown="span">
  ![Detailansicht Datei 3](img/file_details_3_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 6: Detailansicht einer fehlerhaften Datei</figcaption></center>
</figure>

Im Reiter Metadaten werden alle Eigenschaften die von Borg extrahiert wurden gruppiert nach Kategorie aufgelistet.

<figure markdown="span">
  ![Datei Metadaten](img/file_metadata.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 7: Metadatenansicht einer Videodatei</figcaption></center>
</figure>
