# Job Details​Copy link
- Get all job details, including additional information such as: application options / links, employer reviews and estimated salaries for similar jobs.

## Query Parameters
- job_id
  - Type:string
  - required
  - Example
  - Job Id of the job for which to get details.
  - Batching of up to 20 Job Ids is supported by separating multiple Job Ids by comma (,).
  - Note that each Job Id in a batch request is counted as a request for quota calculation.

countryCopy link to country
Type:string
default: 
"us"
Example
Country code of the country from which to return job posting.

Allowed values: See https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2

languageCopy link to language
Type:string
Example
Language code in which to return job postings.
Leave empty to use the primary language in the specified country (country parameter).

Allowed values: See https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes

fieldsCopy link to fields
Type:string
Example
A comma separated list of job fields to include in the response (field projection).
By default all fields are returned.

Responses

200
Successful Response

application/json
Request Example forget/job-details

