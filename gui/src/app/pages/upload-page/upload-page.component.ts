import { AfterViewInit, Component, ViewChild, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { FileUpload } from '../../services/file-analysis.service';
import { UploadService } from '../../services/upload.service';
import { FileSizePipe } from '../../shared/file-size.pipe';

@Component({
    selector: 'app-upload-page',
    templateUrl: './upload-page.component.html',
    styleUrls: ['./upload-page.component.scss'],
    imports: [
        FileSizePipe,
        MatButtonModule,
        MatIconModule,
        MatPaginatorModule,
        MatProgressBarModule,
        MatProgressSpinnerModule,
        MatSortModule,
        MatTableModule,
    ]
})
export class UploadPageComponent implements AfterViewInit {
  private upload = inject(UploadService);

  dataSource: MatTableDataSource<FileUpload>;
  displayedColumns: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor() {
    this.dataSource = new MatTableDataSource<FileUpload>();
    this.displayedColumns = [
      'path',
      'filename',
      'fileSize',
      'uploadProgress',
      'verificationProgress',
    ];
    this.upload
      .getAll()
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
      const fileUpload = this.upload.add(file.name, 'Einzeldatei', file.size);
      this.upload.upload(file, fileUpload);
    }
  }

  addFolder(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length > 1) {
      for (let fileIndex = 0; fileIndex < files.length; ++fileIndex) {
        const file = files[fileIndex];
        const fileUpload = this.upload.add(
          file.name,
          file.webkitRelativePath.replace(new RegExp('/' + file.name + '$'), ''),
          file.size,
        );
        this.upload.upload(file, fileUpload);
      }
    }
  }
}
