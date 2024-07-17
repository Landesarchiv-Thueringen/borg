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

import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { FileOverviewComponent } from 'src/app/file-overview/file-overview.component';
import { formatFileSize } from 'src/app/utility/formatting/file-size.pipe';
import { FileFeaturePipe } from 'src/app/utility/localization/file-attribut-de.pipe';
import { Feature, FileAnalysisService, FileResult, OverviewFeature } from '../file-analysis.service';
import { StatusIcons, StatusIconsService } from '../status-icons.service';

type FileOverview = {
  [key in OverviewFeature]?: FileFeature;
} & {
  id: FileFeature;
  icons: StatusIcons;
};

interface FileFeature {
  value: string;
  confidence?: number;
  feature?: Feature;
}

@Component({
  selector: 'app-file-analysis-table',
  templateUrl: './file-analysis-table.component.html',
  styleUrls: ['./file-analysis-table.component.scss'],
  standalone: true,
  imports: [MatTableModule, MatIconModule, MatButtonModule, MatPaginatorModule, FileFeaturePipe, CommonModule],
})
export class FileAnalysisTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileOverview>;
  generatedTableColumnList: string[];
  tableColumnList: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;

  constructor(
    private dialog: MatDialog,
    private fileAnalysisService: FileAnalysisService,
    private statusIcons: StatusIconsService,
  ) {
    this.dataSource = new MatTableDataSource<FileOverview>([]);
    this.tableColumnList = [];
    this.generatedTableColumnList = ['fileName', 'relativePath', 'fileSize'];
    this.fileAnalysisService.getFileResults().subscribe({
      // error can't occur --> no error handling
      next: (fileInfos: FileResult[]) => {
        this.processFileInformation(fileInfos);
      },
    });
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
  }

  processFileInformation(fileInfos: FileResult[]): void {
    const featureKeys: string[] = ['fileName', 'relativePath', 'fileSize'];
    const data: FileOverview[] = [];
    for (let fileInfo of fileInfos) {
      let fileOverview: FileOverview = {
        id: { value: fileInfo.id },
        icons: this.statusIcons.getIcons(fileInfo),
      };
      fileOverview['fileName'] = { value: fileInfo.fileName };
      fileOverview['relativePath'] = { value: fileInfo.relativePath ?? '' };
      fileOverview['fileSize'] = {
        value: formatFileSize(fileInfo.fileSize),
      };
      for (let featureKey in fileInfo.toolResults.summary) {
        if (this.fileAnalysisService.isOverviewFeature(featureKey) && featureKey !== 'valid') {
          featureKeys.push(featureKey);
          fileOverview[featureKey] = {
            value: fileInfo.toolResults.summary[featureKey].values[0].value,
            confidence: fileInfo.toolResults.summary[featureKey].values[0].score,
            feature: fileInfo.toolResults.summary[featureKey],
          };
        }
      }
      data.push(fileOverview);
    }
    this.dataSource.data = data;
    const sortedFeatures = this.fileAnalysisService.sortFeatures(featureKeys);
    this.generatedTableColumnList = sortedFeatures;
    this.tableColumnList = sortedFeatures.concat(['status']);
  }

  openDetails(fileOverview: FileOverview): void {
    const id = fileOverview['id']?.value;
    const fileResult = this.fileAnalysisService.getFileResult(id);
    if (fileResult) {
      this.dialog.open(FileOverviewComponent, {
        data: {
          fileResult: fileResult,
        },
        autoFocus: false,
        maxWidth: '80vw',
      });
    } else {
      console.error('file result not found');
    }
  }

  clearToolResults(): void {
    this.fileAnalysisService.clearFileResults();
  }

  exportResults(): void {
    const a = document.createElement('a');
    document.body.appendChild(a);
    a.download = 'borg-results.json';
    a.href =
      'data:text/json;charset=utf-8,' +
      encodeURIComponent(JSON.stringify(this.fileAnalysisService.fileResults, null, 2));
    a.click();
    document.body.removeChild(a);
  }
}
