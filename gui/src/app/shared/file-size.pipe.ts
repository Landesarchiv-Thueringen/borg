import { formatNumber } from '@angular/common';
import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'fileSize',
  standalone: true,
})
export class FileSizePipe implements PipeTransform {
  transform(value: number): string {
    return formatFileSize(value);
  }
}

const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB'];

export function formatFileSize(value: number): string {
  let exp = Math.floor(Math.log10(value));
  exp = exp - (exp % 3);
  const sizeMB = value / Math.pow(10, exp);
  let unitIndex = Math.floor(exp / 3);
  // gigantic file, something is wrong
  if (unitIndex > units.length) {
    console.error('unrealistic large file');
    return value + units[0];
  }

  return formatNumber(sizeMB, 'de', '1.0-2') + ' ' + units[unitIndex];
}
