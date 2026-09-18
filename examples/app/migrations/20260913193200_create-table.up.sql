-- projects definition

CREATE TABLE IF NOT EXISTS projects (
	id TEXT NOT NULL, disp_name TEXT NOT NULL,
	CONSTRAINT projects_pk PRIMARY KEY (id)
);

-- activities definition

CREATE TABLE IF NOT EXISTS activities (
	id TEXT NOT NULL,
	project_id TEXT NOT NULL,
	disp_name TEXT NOT NULL,
	duration INTEGER NOT NULL,
	CONSTRAINT activities_pk PRIMARY KEY (id),
	CONSTRAINT activities_projects_FK FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- dependencies definition

CREATE TABLE IF NOT EXISTS dependencies (
	id TEXT NOT NULL,
	project_id TEXT NOT NULL,
	relationship TEXT NOT NULL,
	predecessor_activity_id TEXT NOT NULL,
	successor_activity_id TEXT NOT NULL,
	CONSTRAINT dependencies_pk PRIMARY KEY (id),
	CONSTRAINT dependencies_activities_predecessor_FK FOREIGN KEY (predecessor_activity_id) REFERENCES activities(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT dependencies_activities_successor_FK FOREIGN KEY (successor_activity_id) REFERENCES activities(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT dependencies_projects_FK FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE
);