import { CommonModule } from '@angular/common';
import { HttpModule } from '@angular/http';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import { DashboardComponent } from './containers/dashboard/dashboard.component';

import { SchemaComponent } from './components/schema/schema.component';
import { ExecutionHistoryComponent } from './components/execution-history/execution-history.component';
import { ExecutionTemplatesComponent } from './components/execution-templates/execution-templates.component';
import { QueryEditorComponent } from './components/query-editor/query-editor.component';

import { FilterPipe } from './pipes/filter.pipe';

import { SchemaReducer } from './reducers/schema.reducer';
import { SchemaEffects } from './effects/schema.effects';
import { SchemaService } from './services/schema.service';

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
    FilterPipe,
    DashboardComponent,
    SchemaComponent,
    ExecutionHistoryComponent,
    ExecutionTemplatesComponent,
    QueryEditorComponent,
  ],
  imports: [
    CommonModule,
    HttpModule,
    RouterModule.forRoot(ROUTES),
    FormsModule,
    StoreModule.forFeature('dashboard', {
      schema: SchemaReducer,
    }),
    EffectsModule.forFeature([SchemaEffects]),

  ],
  providers: [
    SchemaService,
  ],
})
export class DashboardModule {}
