import { Component, input, OnInit } from '@angular/core';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTableModule } from '@angular/material/table';
import { FileAnalysis } from '../results';

interface Feature {
  category: string;
  key: string;
  value: string | number | boolean;
  supportingTools: string[];
}

@Component({
  selector: 'app-metadata-table',
  imports: [MatExpansionModule, MatTableModule],
  templateUrl: './metadata-table.component.html',
  styleUrl: './metadata-table.component.scss',
})
export class MetadataTableComponent implements OnInit {
  fileAnalysis = input.required<FileAnalysis>();
  displayedColumns: string[] = ['category', 'key', 'value', 'tools'];
  features: Feature[] = [];

  ngOnInit(): void {
    if (this.fileAnalysis().featureSets.length > 0) {
      for (let key in this.fileAnalysis().featureSets[0].features) {
        const parts = key.split(':');
        if (parts.length !== 2) {
          console.error('Could not extract category and attribute key from: ' + key);
          continue;
        }
        const categoryKey = parts[0];
        const featureKey = parts[1];
        this.features.push({
          category: categoryKey,
          key: featureKey,
          value: this.fileAnalysis().featureSets[0].features[key].value,
          supportingTools: this.fileAnalysis().featureSets[0].features[key].supportingTools,
        });
      }
      console.log(this.features);
    }
  }
}
