import * as fs from "node:fs";
import * as path from "node:path";

function expect(condition, location, message) {
	if (!condition) throw new Error(`${location}: ${message}`);
}

const string = (value, location) => expect(typeof value === "string", location, "expected a string");
const text = (value, location) => {
	string(value, location);
	expect(value.trim().length > 0, location, "expected a non-empty string");
};
const integer = (value, location) => expect(Number.isSafeInteger(value), location, "expected an integer");
const positiveInteger = (value, location) => {
	integer(value, location);
	expect(value > 0, location, "expected a positive integer");
};
const literal = (expected) => (value, location) => expect(value === expected, location, `expected ${JSON.stringify(expected)}`);
const array = (item) => (value, location) => {
	expect(Array.isArray(value), location, "expected an array");
	value.forEach((entry, index) => item(entry, `${location}[${index}]`));
};
const object = (required, optional = {}) => (value, location) => {
	expect(value !== null && typeof value === "object" && !Array.isArray(value), location, "expected an object");
	for (const [key, validate] of Object.entries(required)) validate(value[key], `${location}.${key}`);
	for (const [key, validate] of Object.entries(optional)) {
		if (Object.hasOwn(value, key)) validate(value[key], `${location}.${key}`);
	}
	for (const key of Object.keys(value)) {
		expect(Object.hasOwn(required, key) || Object.hasOwn(optional, key), `${location}.${key}`, "unexpected field; use the documented manifest shape");
	}
};

const strings = array(string);
const section = object({ kind: text, title: text, summary: text }, { points: strings });
const capture = object({ path: text, startLine: positiveInteger, endLine: positiveInteger, snippet: text }, { label: string, reason: string });
const breakpoint = object({
	id: text, order: positiveInteger, title: text, paths: strings, reason: text, sketch: text, captures: array(capture),
}, { dependsOn: strings, downstream: strings, sections: array(section) });
const apiExample = object({ title: text, method: text, path: text }, {
	summary: string,
	requestBody: () => {}, // Bodies are arbitrary JSON values, including null.
	responseStatus: (value, location) => {
		integer(value, location);
		expect(value >= 100 && value <= 599, location, "expected an HTTP status from 100 to 599");
	},
	responseBody: () => {},
});
const report = object({ id: text, title: text, status: text }, {
	coverage: strings,
	findings: array(object({ title: text, evidence: text, impact: text }, { threadId: string, suggestedComment: string })),
	nonFindings: strings,
});
const manifest = object({
	mode: literal("learn"),
	version: literal(1),
	goal: text,
	scope: text,
	overview: object({ summary: text }, { sections: array(section), api: array(apiExample) }),
	breakpoints: array(breakpoint),
}, {
	skillReports: array((value, location) => {
		expect(value !== null && typeof value === "object" && !Array.isArray(value), location,
			"expected a report object; put prose notes in overview.sections or leave skillReports empty");
		report(value, location);
	}),
});

function validateManifest(value) {
	manifest(value, "manifest");
	expect(value.breakpoints.length > 0, "manifest.breakpoints", "expected at least one source step");
	const ids = new Set();
	value.breakpoints.forEach((step, index) => {
		const location = `manifest.breakpoints[${index}]`;
		expect(/^[a-zA-Z0-9_-]+$/.test(step.id), `${location}.id`, "use only letters, digits, underscores, or hyphens");
		expect(!ids.has(step.id), `${location}.id`, `duplicate ID ${JSON.stringify(step.id)}`);
		ids.add(step.id);
		expect(step.captures.length > 0, `${location}.captures`, "expected at least one source capture");
		step.captures.forEach((entry, captureIndex) => {
			expect(entry.endLine >= entry.startLine, `${location}.captures[${captureIndex}].endLine`, "must be at least startLine");
		});
	});
	value.breakpoints.forEach((step, index) => {
		for (const field of ["dependsOn", "downstream"]) {
			(step[field] ?? []).forEach((id, referenceIndex) => {
				expect(ids.has(id), `manifest.breakpoints[${index}].${field}[${referenceIndex}]`, `unknown step ID ${JSON.stringify(id)}`);
			});
		}
	});
}

if (process.argv.length !== 3) {
	console.error("Usage: node validate-manifest.mjs <artifact-directory-or-manifest.json>");
	process.exitCode = 2;
} else {
	try {
		const target = path.resolve(process.argv[2]);
		const filename = fs.statSync(target).isDirectory() ? path.join(target, "manifest.json") : target;
		validateManifest(JSON.parse(fs.readFileSync(filename, "utf8")));
		console.log(`Valid learn manifest: ${filename}`);
	} catch (error) {
		console.error(`Invalid learn manifest: ${error.message}`);
		process.exitCode = 1;
	}
}
