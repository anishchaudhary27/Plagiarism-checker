import {
  Heading,
  Input,
  FormControl,
  FormLabel,
  FormHelperText,
  Button,
  Text,
  CircularProgress,
  useToast,
} from "@chakra-ui/react";
import styles from "./style.module.css";
import {useNavigate} from 'react-router-dom'
import { FaGithub, FaUpload } from "react-icons/fa";
import { useState } from "react";

export default function Home() {
  const [uploading, setUploading] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const toast = useToast();
  const navigate = useNavigate()
  const handleFileChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
    const file = e.target.files![0];
    if (!file) return;
    setFile(file);
  };
  const handleUpload = () => {
    if (file) {
      setUploading(true);
      const formData = new FormData();
      formData.append("file", file);
      fetch("/api/submit", {
        method: "POST",
        body: formData,
      })
        .then((res) => res.json())
        .then((res) => {
          setUploading(false);
          navigate(`/job?id=${res.jobId}`)
        })
        .catch(err => {
          console.error(err)
          setUploading(false);
          toast({
            title: "Error uploading file",
            status: "error",
            duration: 5000,
            isClosable: true,
          });
        })
    } else {
      toast({
        title: "File not selected",
        description: "Please select a file to upload",
        status: "error",
        duration: 5000,
        isClosable: true,
      });
    }
  };
  return (
    <div className={styles.root}>
      <Heading color="teal" size="3xl">
        Plagiarism checker
      </Heading>
      <Text fontSize="lg" marginTop="10">
      An online plagiarism checker service that checks PDFs for plagiarism against articleson wikipedia(en). 
      Tech stack includes Go, PostgreSQL, RabbitMQ, Google Custom Search API, React, Typescript.
      </Text>
      <FormControl marginTop="20">
        <FormLabel as="legend">Choose PDF</FormLabel>
        <Input type="file" accept=".pdf" onChange={handleFileChange} />
        <FormHelperText>
          Choose PDF you want to check for plagiarism
        </FormHelperText>
      </FormControl>
      {uploading ? (
        <CircularProgress isIndeterminate color="teal" marginTop="10" />
      ) : (
        <Button
          marginTop="10"
          colorScheme="teal"
          leftIcon={<FaUpload />}
          onClick={handleUpload}
        >
          Upload
        </Button>
      )}
      <Button
        marginTop="10"
        leftIcon={<FaGithub />}
        onClick={() =>
          window.open("https://github.com/anishchaudhary27/Plagiarism-checker")
        }
      >
        Visit github repository
      </Button>
    </div>
  );
}
