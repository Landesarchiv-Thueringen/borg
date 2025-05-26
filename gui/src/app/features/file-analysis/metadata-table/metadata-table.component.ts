import { Component, computed, input, OnInit, signal, WritableSignal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { Localization, LocalizationPipe } from '../pipes/localization.pipe';
import { FileAnalysis } from '../results';

interface Feature {
  category: string;
  key: string;
  value: string | number | boolean;
  supportingTools: string[];
}

interface FeatureFilter {
  category: string | null;
}

@Component({
  selector: 'app-metadata-table',
  imports: [
    MatExpansionModule,
    MatTableModule,
    MatFormFieldModule,
    MatSelectModule,
    ReactiveFormsModule,
    LocalizationPipe,
  ],
  templateUrl: './metadata-table.component.html',
  styleUrl: './metadata-table.component.scss',
})
export class MetadataTableComponent implements OnInit {
  readonly fileAnalysis = input.required<FileAnalysis>();
  readonly localization = input.required<Localization | undefined>();
  displayedColumns: string[] = ['category', 'key', 'value', 'tools'];
  features: Feature[] = [];
  categories: string[] = [];
  filterForm = new FormGroup({
    category: new FormControl(''),
  });
  filter: WritableSignal<FeatureFilter | undefined> = signal(undefined);
  filteredFeatures = computed(() => this.filterFeatures(this.features, this.filter()));

  constructor() {
    this.filterForm.controls.category.valueChanges.subscribe((value) => {
      this.filter.set({
        category: value,
      });
    });
  }

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
          key: key,
          value: this.fileAnalysis().featureSets[0].features[key].value,
          supportingTools: this.fileAnalysis().featureSets[0].features[key].supportingTools,
        });
        if (!this.categories.includes(categoryKey)) {
          this.categories.push(categoryKey);
        }
      }
    }
  }

  filterFeatures(features: Feature[], filter: FeatureFilter | undefined): Feature[] {
    const filteredFeatures: Feature[] = [];
    for (let f of features) {
      if (!filter) {
        filteredFeatures.push(f);
      } else if (this.featureFilterApplies(f, filter)) {
        filteredFeatures.push(f);
      }
    }
    return filteredFeatures;
  }

  featureFilterApplies(feature: Feature, filter: FeatureFilter) {
    if (filter.category && feature.category !== filter.category) {
      return false;
    }
    return true;
  }

  setCategoryFilter(event: any) {
    console.log(event);
  }
}
