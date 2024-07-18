import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, Input, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { isOverviewFeature, OverviewFeature, sortFeatures } from '../file-feature';
import { FileOverviewComponent } from '../file-overview/file-overview.component';
import { FileFeaturePipe } from '../pipes/file-attribut-de.pipe';
import { formatFileSize } from '../pipes/file-size.pipe';
import { Feature, FileResult } from '../results';
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
  imports: [
    CommonModule,
    FileFeaturePipe,
    MatButtonModule,
    MatIconModule,
    MatPaginatorModule,
    MatSortModule,
    MatTableModule,
  ],
})
export class FileAnalysisTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileOverview>;
  generatedTableColumnList: string[];
  tableColumnList: string[];

  private _results?: FileResult[];
  @Input()
  get results(): FileResult[] | undefined {
    return this._results;
  }
  set results(value: FileResult[] | undefined) {
    this._results = value;
    this.processFileInformation(value ?? []);
  }
  @Input() getResult!: (id: string) => Promise<FileResult | undefined>;

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor(
    private dialog: MatDialog,
    private statusIcons: StatusIconsService,
  ) {
    this.dataSource = new MatTableDataSource<FileOverview>([]);
    this.tableColumnList = [];
    this.generatedTableColumnList = ['fileName', 'relativePath', 'fileSize'];
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
    this.dataSource.sortingDataAccessor = (data, sortHeaderId) => {
      switch (sortHeaderId) {
        case 'status':
          return Object.entries(data.icons)
            .filter(([_, value]) => value)
            .map(([key, _]) => key)
            .join();
        default:
          const feature = data[sortHeaderId as keyof FileOverview] as FileFeature;
          return feature?.value;
      }
    };
  }

  private processFileInformation(fileInfos: FileResult[]): void {
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
        if (isOverviewFeature(featureKey) && featureKey !== 'valid') {
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
    const sortedFeatures = sortFeatures(featureKeys);
    this.generatedTableColumnList = sortedFeatures;
    this.tableColumnList = sortedFeatures.concat(['status']);
  }

  async openDetails(fileOverview: FileOverview): Promise<void> {
    const id = fileOverview['id']?.value;
    const fileResult = await this.getResult(id);
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
}
