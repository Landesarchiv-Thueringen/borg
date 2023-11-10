// angular
import { Component, Inject } from '@angular/core';

// material
import { MAT_DIALOG_DATA } from '@angular/material/dialog';

// project
import { FileOverview } from '../file-analysis/file-analysis-table/file-analysis-table.component';

export interface DialogData {
  fileOverview: FileOverview;
}

@Component({
  selector: 'app-file-overview',
  templateUrl: './file-overview.component.html',
  styleUrls: ['./file-overview.component.scss']
})
export class FileOverviewComponent {
  readonly fileOverview: FileOverview;
  constructor(@Inject(MAT_DIALOG_DATA) private data: DialogData) {
    this.fileOverview = data.fileOverview;
  }
}
