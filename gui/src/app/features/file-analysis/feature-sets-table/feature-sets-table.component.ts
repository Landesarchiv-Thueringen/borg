import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { FileFeaturePipe } from '../pipes/file-feature.pipe';
import { FeatureSet } from '../results';

interface DialogData {
  featureSets: FeatureSet[];
}

interface FeatureValue {
  key: string;
  value: string | number | boolean;
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
  imports: [CommonModule, MatDialogModule, MatButtonModule, FileFeaturePipe],
})
export class FeatureSetsTableComponent {
  private data = inject<DialogData>(MAT_DIALOG_DATA);
  json: string[] = [];
  ms: Mockup[] = [];
  constructor() {
    for (let f of this.data.featureSets) {
      const fvs: FeatureValue[] = [];
      for (let key in f.features) {
        fvs.push({
          key: key,
          value: f.features[key],
        });
      }
      this.ms.push({
        values: fvs,
        score: f.score,
        tools: f.supportingTools.join(', '),
      });
      this.json.push(JSON.stringify(f));
    }
  }
}
