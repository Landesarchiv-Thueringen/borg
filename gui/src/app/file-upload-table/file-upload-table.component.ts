// angular
import { AfterViewInit, Component, ViewChild } from '@angular/core';

// material
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';

export interface FileUpload {
  fileName: string;
  relativePath?: string;
  fileSize: number;
  uploadProgress?: string;
}

@Component({
  selector: 'app-file-upload-table',
  templateUrl: './file-upload-table.component.html',
  styleUrls: ['./file-upload-table.component.scss'],
})
export class FileUploadTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileUpload>;
  displayedColumns: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor() {
    this.dataSource = new MatTableDataSource<FileUpload>();
    this.displayedColumns = [
      'fileName',
      'relativePath',
      'fileSize',
      'uploadProgress',
    ];
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
      console.log(file);
      const fileUpload: FileUpload = {
        fileName: file.name,
        fileSize: file.size,
      }
      const data = this.dataSource.data
      data.push(fileUpload)
      this.dataSource.data = data
    }
  }

  addFolder(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length > 1) {
      const data = this.dataSource.data
      for (let fileIndex = 0; fileIndex < files.length; ++fileIndex) {
        const file = files[fileIndex];
        console.log(file);
        const fileUpload: FileUpload = {
          fileName: file.name,
          // remove file name from path
          relativePath: file.webkitRelativePath.replace(new RegExp(file.name + '$'), ''),
          fileSize: file.size,
        }
        data.push(fileUpload)
      }
      this.dataSource.data = data
    }
  }
}
