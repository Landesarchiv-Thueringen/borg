import { DecimalPipe } from '@angular/common';
import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'fileSize' })
export class FileSizePipe implements PipeTransform {
  readonly units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB'];

  constructor(private decimalPipe: DecimalPipe) {}

  transform(value: number): string {
    let exp = Math.floor(Math.log10(value));
    exp = exp - (exp % 3);
    const sizeMB = value / Math.pow(10, exp);
    let unitIndex = Math.floor(exp / 3);
    // gigantic file, something is wrong
    if (unitIndex > this.units.length) {
      console.error('unrealistic large file');
      return value + this.units[0];
    }
    return this.decimalPipe.transform(sizeMB, '1.0-2') + ' ' + this.units[unitIndex];
  }
}
