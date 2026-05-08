import { z } from "zod";

const commitSchema = z.object({
  sha: z.string(),
  html_url: z.string().url()
});

export async function fetchLatestRepoCommit() {
  const response = await fetch(
    "https://api.github.com/repos/baditaflorin/localhuman-mail/commits/main",
    {
      headers: {
        Accept: "application/vnd.github+json"
      }
    }
  );
  if (!response.ok) {
    throw new Error("Latest commit unavailable");
  }
  const commit = commitSchema.parse(await response.json());
  return {
    sha: commit.sha.slice(0, 12),
    url: commit.html_url
  };
}
