import { PercentPipe } from '@angular/common';
import { Component, input, OnInit } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { FeatureValuePipe } from '../pipes/feature-value.pipe';
import { FeatureValue, FileAnalysis } from '../results';

interface FormatRow {
  puid: FeatureValue | undefined;
  mimeType: FeatureValue | undefined;
  formatVersion: FeatureValue | undefined;
  valid: FeatureValue | undefined;
  score: number;
  tools: string[];
  errors: boolean;
}

@Component({
  selector: 'app-file-format',
  imports: [MatTableModule, PercentPipe, FeatureValuePipe, MatIconModule],
  templateUrl: './file-format.component.html',
  styleUrl: './file-format.component.scss',
})
export class FileFormatComponent implements OnInit {
  readonly fileAnalysis = input.required<FileAnalysis>();
  displayedColumns: string[] = ['puid', 'mimeType', 'formatVersion', 'valid', 'tools', 'score'];
  rows: FormatRow[] = [];

  ngOnInit(): void {
    if (this.fileAnalysis().featureSets.length > 0) {
      let setIndex = 0;
      this.rows = this.fileAnalysis().featureSets.map((set) => {
        setIndex += 1;
        return {
          puid: set.features['format:puid'],
          mimeType: set.features['format:mimeType'],
          formatVersion: set.features['format:version'],
          valid: set.features['format:valid'],
          score: set.score,
          tools: set.supportingTools,
          errors: this.fileAnalysis().toolResults.some(
            (tr) => set.supportingTools.includes(tr.id) && tr.error,
          ),
        };
      });
    }
  }
}
