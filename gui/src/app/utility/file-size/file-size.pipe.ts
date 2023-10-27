import { Pipe, PipeTransform } from '@angular/core';

@Pipe({name: 'fileSize'})
export class FileSizePipe implements PipeTransform {
  transform(value: number): string {
    const sizeMB = value / 1000000;
    return sizeMB.toString() + ' MB';
  }
}