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

import { Injectable } from '@angular/core';
import { MatPaginatorIntl } from '@angular/material/paginator';

@Injectable({
  providedIn: 'root',
})
export class PaginatorDeService extends MatPaginatorIntl {
  override itemsPerPageLabel = 'Einträge pro Seite';
  override nextPageLabel = 'nächste Seite';
  override previousPageLabel = 'vorherige Seite';
  override firstPageLabel = 'erste Seite';
  override lastPageLabel = 'letzte Seite';
  override getRangeLabel = (page: number, pageSize: number, length: number): string => {
    const preposition = 'von';
    if (length === 0 || pageSize === 0) {
      return '';
    }
    length = Math.max(length, 0);
    const startIndex = page * pageSize;
    // If the start index exceeds the list length, do not try and fix the end index to the end.
    const endIndex = startIndex < length ? Math.min(startIndex + pageSize, length) : startIndex + pageSize;
    return startIndex + 1 + ' - ' + endIndex + ' ' + preposition + ' ' + length;
  };
}
