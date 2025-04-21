import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { PrettyPrintJsonPipe } from '../pipes/pretty-print-json.pipe';
import { FeatureSet } from '../results';

interface DialogData {
  faetureSets: FeatureSet[];
}

interface Mockup {
  values: string;
  tools: string;
  score: number;
}

@Component({
  selector: 'app-feature-sets-table',
  templateUrl: './feature-sets-table.component.html',
  styleUrls: ['./feature-sets-table.component.scss'],
  imports: [CommonModule, MatDialogModule, PrettyPrintJsonPipe],
})
export class FeatureSetsTableComponent {
  private data = inject<DialogData>(MAT_DIALOG_DATA);
  json: string[] = [];
  ms: Mockup[] = [];
  constructor() {
    for (let f of this.data.faetureSets) {
      this.ms.push({
        values: JSON.stringify(f.features),
        score: f.score,
        tools: f.supportingTools.join(', '),
      });
      this.json.push(JSON.stringify(f));
    }
  }
}
