// angular
import { NgModule, LOCALE_ID } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DecimalPipe } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { registerLocaleData } from '@angular/common';
import localeDe from '@angular/common/locales/de';
registerLocaleData(localeDe);

// material
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog'; 
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu'; 
import {
  MatPaginatorModule,
  MatPaginatorIntl,
} from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'; 
import { MatTableModule } from '@angular/material/table';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';

// project
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MainNavigationComponent } from './main-navigation/main-navigation.component';
import { FileFeaturePipe } from './utility/localization/file-attribut-de.pipe';
import { FileSizePipe } from './utility/formatting/file-size.pipe';
import { FileUploadTableComponent } from './file-upload-table/file-upload-table.component';
import { PaginatorDeService } from './utility/localization/paginator-de.service';
import { FileAnalysisTableComponent } from './file-analysis/file-analysis-table/file-analysis-table.component';
import { FileOverviewComponent } from './file-overview/file-overview.component';
import { ToolOutputComponent } from './tool-output/tool-output.component';
import { CsvTablePipe } from './utility/formatting/csv-table.pipe';
import { PrettyPrintJsonPipe } from './utility/formatting/pretty-print-json.pipe';

@NgModule({
  declarations: [
    AppComponent,
    MainNavigationComponent,
    FileFeaturePipe,
    FileSizePipe,
    FileUploadTableComponent,
    FileAnalysisTableComponent,
    FileOverviewComponent,
    ToolOutputComponent,
    CsvTablePipe,
    PrettyPrintJsonPipe,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    MatMenuModule,
    MatPaginatorModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatSortModule,
    MatTableModule,
    MatToolbarModule,
    MatTooltipModule,
  ],
  providers: [
    { provide: LOCALE_ID, useValue: 'de' },
    { provide: MatPaginatorIntl, useClass: PaginatorDeService },
    DecimalPipe,
    FileFeaturePipe,
    FileSizePipe,
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
