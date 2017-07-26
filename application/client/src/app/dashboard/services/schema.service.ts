import 'rxjs/add/operator/map';
import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { SchemaTable } from '../models/schema.model';

import { environment } from '../../../environments/environment';

@Injectable()
export class SchemaService {
  constructor(private http: Http) {}

  get(): Observable<SchemaTable[]> {
    return this.http
    .get(`${environment.apiPath}/schema`)
    .map(res => this._transformResponse(res.json()));
  }

  _transformResponse(resp: any[]): SchemaTable[] {
    const groups = resp.reduce(function(acc, x) {
      const elem = {
        tableName: x.table_name,
        columnName: x.column_name,
        dataType: x.data_type
      };
      const tableName = elem.tableName;
      (acc[tableName] = acc[tableName] || []).push(elem);
      return acc;
    }, {});

    return Object.keys(groups).map((t) => ({tableName: t, columns: groups[t]}));
  }
}
