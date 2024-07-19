import { Component } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { Router } from '@angular/router';
import { FileAnalysisTableComponent } from '../../features/file-analysis/file-analysis-table/file-analysis-table.component';
import { FileResult } from '../../features/file-analysis/results';
import { ResultsService } from '../../services/results.service';

@Component({
  selector: 'app-results-page',
  standalone: true,
  imports: [FileAnalysisTableComponent, MatIconModule, MatButtonModule],
  templateUrl: './results-page.component.html',
  styleUrl: './results-page.component.scss',
})
export class ResultsPageComponent {
  results?: FileResult[];
  getResult = (id: string) => this.resultsService.get(id);

  constructor(
    private router: Router,
    private resultsService: ResultsService,
  ) {
    this.resultsService
      .getAll()
      .pipe(takeUntilDestroyed())
      .subscribe((results) => (this.results = results));
  }

  clearToolResults(): void {
    this.resultsService.clear();
    this.router.navigate(['auswahl']);
  }

  exportResults(): void {
    const a = document.createElement('a');
    document.body.appendChild(a);
    a.download = 'borg-results.json';
    a.href =
      'data:text/json;charset=utf-8,' +
      encodeURIComponent(JSON.stringify(this.resultsService.fileResults, null, 2));
    a.click();
    document.body.removeChild(a);
  }
}
