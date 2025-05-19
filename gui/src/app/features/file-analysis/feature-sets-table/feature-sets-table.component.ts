import { CommonModule, PercentPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { BreakOpportunitiesPipe } from '../pipes/break-opportunities.pipe';
import { FileFeaturePipe } from '../pipes/file-feature.pipe';
import { FeatureSet } from '../results';

interface DialogData {
  featureSets: FeatureSet[];
}

interface FeatureValue {
  key: string;
  value: string | number | boolean;
}

interface Row {
  puid: string | undefined;
  mimeType: string | undefined;
  version: string | undefined;
  valid: boolean | undefined;
  tools: string;
  score: number;
}

interface Mockup {
  values: FeatureValue[];
  tools: string;
  score: number;
}

@Component({
  selector: 'app-feature-sets-table',
  templateUrl: './feature-sets-table.component.html',
  styleUrls: ['./feature-sets-table.component.scss'],
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    FileFeaturePipe,
    MatTableModule,
    PercentPipe,
    BreakOpportunitiesPipe,
    MatIconModule,
  ],
})
export class FeatureSetsTableComponent {
  private data = inject<DialogData>(MAT_DIALOG_DATA);
  displayedColumns: string[] = ['tools', 'puid', 'mimeType', 'version', 'valid', 'score'];
  dataSource: Row[] = [];
  json: string[] = [];
  ms: Mockup[] = [];
  constructor() {
    this.dataSource = this.data.featureSets.map((fs) => {
      return {
        puid: fs.features['format:puid'] ? (fs.features['format:puid'] as string) : undefined,
        mimeType: fs.features['format:mimeType']
          ? (fs.features['format:mimeType'] as string)
          : undefined,
        version: fs.features['format:version']
          ? (fs.features['format:version'] as string)
          : undefined,
        valid:
          fs.features['format:valid'] !== undefined
            ? (fs.features['format:valid'] as boolean)
            : undefined,
        tools: fs.supportingTools.join('\n'),
        score: fs.score,
      };
    });
  }
}
