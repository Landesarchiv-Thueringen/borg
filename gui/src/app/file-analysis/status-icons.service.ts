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
import { FileResult } from './file-analysis.service';

const UNCERTAIN_REQUIRED_FEATURES = ['mimeType', 'puid'];
const UNCERTAIN_CONFIDENCE_THRESHOLD = 0.75;
const VALID_CONFIDENCE_THRESHOLD = 0.75;

export interface StatusIcons {
  uncertain: boolean;
  valid: boolean;
  invalid: boolean;
  error: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class StatusIconsService {
  getIcons(fileResult: FileResult): StatusIcons {
    return {
      uncertain: this.hasUncertainIcon(fileResult),
      valid: this.hasValidIcon(fileResult),
      invalid: this.hasInvalidIcon(fileResult),
      error: this.hasErrorIcon(fileResult),
    };
  }

  private hasUncertainIcon(fileResult: FileResult): boolean {
    for (const key of UNCERTAIN_REQUIRED_FEATURES) {
      const feature = fileResult.toolResults.summary[key];
      if (!feature || feature.values[0].score < UNCERTAIN_CONFIDENCE_THRESHOLD) {
        return true;
      }
    }
    return false;
  }

  private hasValidIcon(fileResult: FileResult): boolean {
    const valid = fileResult.toolResults.summary['valid'];
    return (
      !this.hasUncertainIcon(fileResult) &&
      valid?.values[0].value === 'true' &&
      valid.values[0].score > VALID_CONFIDENCE_THRESHOLD
    );
  }

  private hasInvalidIcon(fileResult: FileResult): boolean {
    const valid = fileResult.toolResults.summary['valid'];
    return (
      !this.hasUncertainIcon(fileResult) &&
      valid?.values[0].value === 'false' &&
      valid.values[0].score > VALID_CONFIDENCE_THRESHOLD
    );
  }

  private hasErrorIcon(fileResult: FileResult): boolean {
    return (
      fileResult.toolResults.fileIdentificationResults?.some((result) => result.error) ||
      fileResult.toolResults.fileValidationResults?.some((result) => result.error) ||
      false
    );
  }
}
