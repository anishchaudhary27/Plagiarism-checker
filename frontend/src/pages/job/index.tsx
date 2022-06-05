import {
  Box,
  CircularProgress,
  Heading,
  useToast,
  Badge,
  Link,
  Button,
} from "@chakra-ui/react";
import { ExternalLinkIcon } from "@chakra-ui/icons";
import { useEffect, useState } from "react";
import styles from "./style.module.css";

interface Page {
  title: string;
  description: string;
  similarity: number;
  link: string;
}

export default function Job() {
  const [loading, setloading] = useState(true);
  const [pages, setPages] = useState<Page[]>([]);
  const toast = useToast();
  useEffect(() => {
    const jobId = new URLSearchParams(window.location.search).get("id");
    let looper: null | NodeJS.Timeout = null;
    if (jobId) {
      looper = setInterval(() => {
        fetch("/api/status/" + jobId)
          .then((res) => res.json())
          .then((res) => {
            if (res.status === 1) {
              if (looper) clearInterval(looper);
              setloading(false);
              setPages(res.result);
            }
            else if(res.status === -1) {
              if (looper) clearInterval(looper);
              setloading(false);
              toast({
                title: "Error",
                description: "Error occured while processing document",
                status: "error",
                duration: 9000,
                isClosable: true,
              });
            }
          })
          .catch((error) => {
            console.error(error);
            toast({
              title: "Error getting job status. Trying again...",
              status: "error",
              duration: 1000,
              isClosable: true,
            });
          });
      }, 5000);
    } else {
      toast({
        title: "Invalid job id",
        description: "Invalid job id provided",
        status: "error",
        duration: 5000,
      });
    }
    return () => {
      if (looper) {
        clearInterval(looper);
      }
    };
  }, []);
  return (
    <div className={styles.root}>
      <Heading color="teal" size="2xl">
        Similar articles
      </Heading>
      {loading ? (
        <div className={styles.analysing}>
          <CircularProgress isIndeterminate color="teal.400" />
          <Heading size="lg" marginTop="5" color="teal.400">
            analyzing document
          </Heading>
        </div>
      ) : (
        <div className={styles.container}>
          {pages.map((page, i) => (
            <Box
              maxW={{ base: "100%", md: "50%" }}
              borderWidth="1px"
              borderRadius="lg"
              overflow="hidden"
              key={i}
              marginTop="2"
              backgroundColor="white"
              padding="5"
            >
              <Box
                mt="1"
                fontWeight="semibold"
                as="h2"
                lineHeight="tight"
                noOfLines={1}
              >
                {page.title}
              </Box>
              <Badge
                borderRadius="full"
                px="2"
                colorScheme="teal"
                marginRight="4"
              >
                Similarity: {page.similarity}%
              </Badge>
              <Link href={page.link} isExternal color="blue.500">
                {page.link} <ExternalLinkIcon mx="2px" />
              </Link>
              <Box
                mt="1"
                fontWeight="regular"
                as="p"
                color="gray.500"
                lineHeight="tight"
                noOfLines={2}
              >
                {page.description}
              </Box>
            </Box>
          ))}
        </div>
      )}
    </div>
  );
}
