import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { FileAnalysisService, FileUpload } from '../file-analysis/file-analysis.service';
import { FileSizePipe } from '../utility/formatting/file-size.pipe';
import { UploadService } from '../utility/upload.service';

@Component({
  selector: 'app-file-upload-table',
  templateUrl: './file-upload-table.component.html',
  styleUrls: ['./file-upload-table.component.scss'],
  standalone: true,
  imports: [
    FileSizePipe,
    MatButtonModule,
    MatIconModule,
    MatPaginatorModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSortModule,
    MatTableModule,
  ],
})
export class FileUploadTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileUpload>;
  displayedColumns: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor(
    private fileAnalysisService: FileAnalysisService,
    private upload: UploadService,
  ) {
    this.dataSource = new MatTableDataSource<FileUpload>();
    this.displayedColumns = ['relativePath', 'fileName', 'fileSize', 'uploadProgress', 'verificationProgress'];
    this.fileAnalysisService
      .getFileUploads()
      .pipe(takeUntilDestroyed())
      .subscribe({
        // error can't occur --> no error handling
        next: (fileUploads: FileUpload[]) => {
          this.dataSource.data = fileUploads;
        },
      });
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
  }

  addFile(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length === 1) {
      const file = files[0];
      const fileUpload = this.fileAnalysisService.addFileUpload(file.name, 'Einzeldatei', file.size);
      this.upload.uploadFile(file, fileUpload);
    }
  }

  addFolder(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length > 1) {
      for (let fileIndex = 0; fileIndex < files.length; ++fileIndex) {
        const file = files[fileIndex];
        const fileUpload = this.fileAnalysisService.addFileUpload(
          file.name,
          file.webkitRelativePath.replace(new RegExp('/' + file.name + '$'), ''),
          file.size,
        );
        this.upload.uploadFile(file, fileUpload);
      }
    }
  }
}
