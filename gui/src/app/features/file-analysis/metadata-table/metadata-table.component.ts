import { Component, input, OnInit } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { CategoryPipe } from '../pipes/category.pipe';
import { FileAnalysis } from '../results';

interface Category {
  id: string;
  features: Feature[];
}

interface Feature {
  key: string;
  label: string | null;
  value: string | number | boolean;
  supportingTools: string[];
}

@Component({
  selector: 'app-metadata-table',
  imports: [MatTableModule, CategoryPipe],
  templateUrl: './metadata-table.component.html',
  styleUrl: './metadata-table.component.scss',
})
export class MetadataTableComponent implements OnInit {
  readonly fileAnalysis = input.required<FileAnalysis>();
  displayedColumns: string[] = ['key', 'value', 'tools'];
  categories: Category[] = [];
  categoryOrder: { [key: string]: number | undefined } = {
    general: 1,
    format: 2,
    av_container: 3,
    audio: 4,
    video: 5,
  };

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
        const feature: Feature = {
          key: featureKey,
          label: this.fileAnalysis().featureSets[0].features[key].label,
          value: this.fileAnalysis().featureSets[0].features[key].value,
          supportingTools: this.fileAnalysis().featureSets[0].features[key].supportingTools,
        };
        const category = this.categories.find((c) => c.id === categoryKey);
        if (!category) {
          this.categories.push({
            id: categoryKey,
            features: [feature],
          });
        } else {
          category.features.push(feature);
        }
      }
      this.sortCategories();
    }
  }

  sortCategories(): void {
    this.categories = this.categories.sort((a, b) => {
      const oa = this.categoryOrder[a.id];
      const ob = this.categoryOrder[b.id];
      if (oa && ob) {
        return oa - ob;
      } else if (oa) {
        return 1;
      } else if (ob) {
        return -1;
      }
      return a.id.localeCompare(b.id);
    });
  }
}
