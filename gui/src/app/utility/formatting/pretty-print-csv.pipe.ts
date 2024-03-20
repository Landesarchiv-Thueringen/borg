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

import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'prettyPrintCsv',
})
export class PrettyPrintCsvPipe implements PipeTransform {
  transform(csvString: string): string {
    const rows = csvString.split('\n').map((row) => row.split(',').map((cell) => cell.trim()));
    const nColumns = Math.max(...rows.map((row) => row.length));
    for (let i = 0; i < nColumns; i++) {
      const columnWidth = Math.max(...rows.map((row) => row[i]?.length ?? 0));
      rows.forEach((row) => {
        if (i < row.length - 1) {
          row[i] = padToLength(row[i] ?? '', columnWidth);
        }
      });
    }
    return rows.map((row) => row.join(', ')).join('\n');
  }
}

function padToLength(s: string, length: number): string {
  return s + ' '.repeat(length - s.length);
}
