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

const labelMap: { [key: string]: string | undefined } = {
  tool: 'Werkzeug',
  fileName: 'Dateiname',
  relativePath: 'Pfad',
  fileSize: 'Dateigröße',
  formatVersion: 'Formatversion',
  mimeType: 'MIME-Type',
  puid: 'PUID',
  valid: 'Valide',
  wellFormed: 'Wohlgeformt',
  encoding: 'Zeichenkodierung',
  error: 'Fehler',
};

@Pipe({ name: 'fileFeature' })
export class FileFeaturePipe implements PipeTransform {
  transform(value: string): string {
    const label = labelMap[value];
    return label ?? value;
  }
}
