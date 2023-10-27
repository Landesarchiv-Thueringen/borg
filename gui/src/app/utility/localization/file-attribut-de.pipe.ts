import { Pipe, PipeTransform } from '@angular/core';

@Pipe({name: 'fileAttribute'})
export class FileAttributePipe implements PipeTransform {
  transform(value: string): string {
    switch (value) {
      case 'fileName': {
        return 'Dateiname';
      }
      case 'relativePath': {
        return 'Pfad';
      }
      case 'fileSize': {
        return 'Dateigröße';
      }
      case 'formatVersion': {
        return 'Formatversion';
      }
      case 'mimeType': {
        return 'MIME-Type';
      }
      case 'puid': {
        return 'PUID';
      }
      case 'valid': {
        return 'Validierung';
      }
      case 'wellFormed': {
        return 'wohlgeformt';
      }
      case 'encoding': {
        return 'Kodierung';
      }
    }
    return value;
  }
}