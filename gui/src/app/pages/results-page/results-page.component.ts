import { Component } from '@angular/core';
import { FileAnalysisTableComponent } from '../../features/file-analysis/file-analysis-table/file-analysis-table.component';

@Component({
  selector: 'app-results-page',
  standalone: true,
  imports: [FileAnalysisTableComponent],
  templateUrl: './results-page.component.html',
  styleUrl: './results-page.component.scss',
})
export class ResultsPageComponent {}
