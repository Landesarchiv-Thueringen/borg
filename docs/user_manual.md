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

Die Detailansicht eines Dateiergebnisses ist in zwei Reiter unterteilt, um eine klare und strukturierte Darstellung der Analyseinformationen zu gewährleisten.

Im ersten Reiter mit dem Titel _Überblick_ erhalten Sie eine allgemeine Zusammenfassung des Analyseergebnisses. Hier wird zunächst der Status der Analyse erläutert – beispielsweise, ob die Verarbeitung erfolgreich war, ob Warnungen aufgetreten sind oder ob es Fehler gab.

Darüber hinaus werden in diesem Reiter alle erkannten Dateiformate aufgelistet, die bei der Analyse der Datei identifiziert wurden. Diese Übersicht dient als Einstiegspunkt, um sich einen schnellen Eindruck vom Inhalt und von der Interpretation der Datei durch die verschiedenen Analysewerkzeuge zu verschaffen.

<figure markdown="span">
  ![Detailansicht Datei 1](img/file_details_1_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 3: Detailansicht einer validen JPEG-Datei</figcaption></center>
</figure>

Im zweiten Reiter _Metadaten_ werden alle Eigenschaften angezeigt, die durch Borg während der Analyse aus der Datei extrahiert wurden. Die Informationen sind thematisch nach Kategorien gruppiert, um eine bessere Übersicht und schnellere Orientierung zu ermöglichen.

Der Fokus liegt hierbei vor allem auf technischen und inhaltlichen Metadaten. Dazu zählen beispielsweise Angaben zum Dateiformat, verwendete Codecs bei Audio- und Videodateien (wie z. B. Video-Codec, Auflösung oder Bildrate) sowie inhaltliche Informationen wie Titel, Autor, Erstellungsdatum oder andere eingebettete Metadaten.

<figure markdown="span">
  ![Datei Metadaten](img/file_metadata.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 4: Metadatenansicht einer Videodatei</figcaption></center>
</figure>

In der Anwendung können Sie sich zu jedem erkannten Dateiformat eine Detailansicht anzeigen lassen. Klicken Sie dazu einfach auf eine Zeile in der Übersichtstabelle der erkannten Formate. Daraufhin öffnet sich ein neuer Dialog, der Ihnen eine tiefere Einsicht in die Zusammensetzung der ermittelten Eigenschaften bietet.

In dieser Detailansicht finden Sie eine strukturierte Tabelle, in der nachvollziehbar aufgeschlüsselt ist, welche Werkzeuge (Tools) zur Analyse eingesetzt wurden, welche spezifischen Eigenschaften sie jeweils extrahiert haben, und wie diese Einzelergebnisse schließlich zum Gesamtergebnis des erkannten Dateiformats zusammengeführt wurden.

Diese transparente Darstellung hilft dabei, die Herkunft und Zuverlässigkeit einzelner Informationen besser nachzuvollziehen und ermöglicht eine fundierte Bewertung der Analyseergebnisse.

<figure markdown="span">
  ![Detailansicht Datei 1](img/format_result_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 5: Detailansicht für ein erkanntes Dateiformat</figcaption></center>
</figure>

Wenn bei der Analyse mehrere mögliche Dateiformate erkannt werden, unterstützt die Bewertung der Ergebnisse dabei, diese einzuordnen und die wahrscheinlichste Interpretation auszuwählen.

<figure markdown="span">
  ![Detailansicht Datei 2](img/file_details_2_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 6: Detailansicht einer Datei mit mehreren erkannten Dateiformaten</figcaption></center>
</figure>

Werden Probleme mit einer Datei festgestellt wird das klar in der Oberfläche kommuniziert.

<figure markdown="span">
  ![Detailansicht Datei 3](img/file_details_3_cut.png){ loading=lazy }
  <br>
  <center><figcaption>Abb. 7: Detailansicht einer fehlerhaften Datei</figcaption></center>
</figure>
