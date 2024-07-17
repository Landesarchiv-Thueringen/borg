/* BorgFormat - File format identification and validation
 * Copyright (C) 2024 Landesarchiv Thüringen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
