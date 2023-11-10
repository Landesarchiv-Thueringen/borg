// angular
import { Component, Inject } from '@angular/core';

// material
import { MAT_DIALOG_DATA } from '@angular/material/dialog';

// project
import { FileResult } from '../file-analysis/file-analysis.service';

export interface DialogData {
  fileResult: FileResult;
}

@Component({
  selector: 'app-file-overview',
  templateUrl: './file-overview.component.html',
  styleUrls: ['./file-overview.component.scss']
})
export class FileOverviewComponent {
  readonly fileResult: FileResult;
  constructor(@Inject(MAT_DIALOG_DATA) private data: DialogData) {
    this.fileResult = data.fileResult;
    console.log(this.fileResult);
  }
}
