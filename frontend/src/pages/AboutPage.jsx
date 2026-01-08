import { useNavigate } from "react-router-dom";

export default function AboutPage() {
  const navigate = useNavigate();
  return (
    <>
      <h1>About Page</h1>
      <p>This was made to learn Backend with Go and Kubernetes</p>
      <button onClick={() => navigate("/")}>Go Back</button>
    </>
  );
}
