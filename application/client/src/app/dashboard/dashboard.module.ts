import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './containers/dashboard/dashboard.component';

import { SchemaComponent } from './components/schema/schema.component';
import { ExecutionHistoryComponent } from './components/execution-history/execution-history.component';
import { ExecutionTemplatesComponent } from './components/execution-templates/execution-templates.component';
import { QueryEditorComponent } from './components/query-editor/query-editor.component';

// services
// import { BatteryService } from './tesla-battery.service';

export const ROUTES: Routes = [{
  path: 'editor',
  children: [
    { path: '', redirectTo: 'schema', pathMatch: 'full' },
    { path: 'schema', component: DashboardComponent, data: { section: 'schema' }},
    { path: 'history', component: DashboardComponent, data: { section: 'history' }},
    { path: 'saved', component: DashboardComponent, data: { section: 'templates' }},
    { path: 'saved/:id', component: DashboardComponent, data: { section: 'templates' }}
  ]
}];

@NgModule({
  declarations: [
    DashboardComponent,
    SchemaComponent,
    ExecutionHistoryComponent,
    ExecutionTemplatesComponent,
    QueryEditorComponent,
  ],
  imports: [
    CommonModule,
    RouterModule.forRoot(ROUTES),
  ],
  providers: [
    // BatteryService
  ],
})
export class DashboardModule {}
