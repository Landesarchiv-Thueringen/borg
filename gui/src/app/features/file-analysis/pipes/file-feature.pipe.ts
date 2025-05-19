import { Pipe, PipeTransform } from '@angular/core';

const labelMap: { [key in string]?: string } = {
  tool: 'Werkzeug',
  filename: 'Dateiname',
  'format:name': 'Formatbezeichnung',
  'format:version': 'Formatversion',
  'format:mimeType': 'MIME-Type',
  'format:puid': 'PUID',
  'format:valid': 'Valide',
  'format:wellFormed': 'Wohlgeformt',
  'text:encoding': 'Zeichenkodierung',
  'format:isText': 'textbasiertes Dateiformat',
  error: 'Fehler',
  status: 'Status',
};

@Pipe({
  name: 'fileFeature',
  standalone: true,
})
export class FileFeaturePipe implements PipeTransform {
  transform(value: string): string {
    const label = labelMap[value];
    return label ?? value;
  }
}
