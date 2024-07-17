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

import { DecimalPipe } from '@angular/common';
import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'fileSize',
  standalone: true,
})
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
